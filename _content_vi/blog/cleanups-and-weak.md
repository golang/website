---
title: "Từ unique tới cleanups và weak: các công cụ cấp thấp mới cho hiệu quả"
date: 2025-03-06
by:
- Michael Knyszek
tags:
- weak
- cleanup
- finalizer
summary: Weak pointer và cơ chế hoàn tất tốt hơn trong Go 1.24.
---

Trong [bài viết năm ngoái](/blog/unique) về package `unique`, chúng tôi đã nhắc
đến một số tính năng mới khi đó đang được xem xét đề xuất, và chúng tôi rất vui được chia sẻ rằng
từ Go 1.24, chúng nay đã sẵn sàng cho mọi nhà phát triển Go.
Những tính năng mới đó là [hàm `runtime.AddCleanup`](https://pkg.go.dev/runtime#AddCleanup),
hàm xếp một hàm khác vào hàng đợi để chạy khi một đối tượng không còn có thể được truy cập,
và [kiểu `weak.Pointer`](https://pkg.go.dev/weak#Pointer),
kiểu trỏ an toàn tới một đối tượng mà không ngăn đối tượng đó bị bộ gom rác thu hồi.
Kết hợp lại, hai tính năng này đủ mạnh để bạn tự xây dựng package `unique`
của riêng mình!
Hãy cùng tìm hiểu điều gì làm cho chúng hữu ích và khi nào nên dùng chúng.

Lưu ý: đây là các tính năng nâng cao của bộ gom rác.
Nếu bạn chưa quen với các khái niệm thu gom rác cơ bản, chúng tôi
khuyến nghị mạnh mẽ bạn đọc phần mở đầu của [hướng dẫn bộ gom rác](/doc/gc-guide#Introduction).

## Cleanups

Nếu bạn từng dùng finalizer, thì khái niệm cleanup sẽ
khá quen thuộc.
Finalizer là một hàm, được gắn với một đối tượng đã cấp phát bằng cách [gọi
`runtime.SetFinalizer`](https://pkg.go.dev/runtime#SetFinalizer), và sau đó
được bộ gom rác gọi vào một thời điểm nào đó sau khi đối tượng không còn truy cập được.
Ở mức khái quát, cleanup cũng hoạt động như vậy.

Hãy xem một ứng dụng sử dụng tệp ánh xạ bộ nhớ và xem
cleanup có thể giúp như thế nào.

```
//go:build unix

type MemoryMappedFile struct {
	data []byte
}

func NewMemoryMappedFile(filename string) (*MemoryMappedFile, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Get the file's info; we need its size.
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	// Extract the file descriptor.
	conn, err := f.SyscallConn()
	if err != nil {
		return nil, err
	}
	var data []byte
	connErr := conn.Control(func(fd uintptr) {
		// Create a memory mapping backed by this file.
		data, err = syscall.Mmap(int(fd), 0, int(fi.Size()), syscall.PROT_READ, syscall.MAP_SHARED)
	})
	if connErr != nil {
		return nil, connErr
	}
	if err != nil {
		return nil, err
	}
	mf := &MemoryMappedFile{data: data}
	cleanup := func(data []byte) {
		syscall.Munmap(data) // ignore error
	}
	runtime.AddCleanup(mf, cleanup, data)
	return mf, nil
}
```

Một tệp ánh xạ bộ nhớ có nội dung được ánh xạ vào bộ nhớ, trong trường hợp này là
dữ liệu nền của một byte slice.
Nhờ phép màu của hệ điều hành, việc đọc và ghi vào byte slice sẽ trực tiếp
truy cập nội dung của tệp.
Với đoạn mã này, chúng ta có thể truyền một `*MemoryMappedFile` đi khắp nơi, và khi nó
không còn được tham chiếu, phép ánh xạ bộ nhớ đã tạo sẽ được dọn dẹp.

Hãy lưu ý rằng `runtime.AddCleanup` nhận ba đối số: địa chỉ của một
biến để gắn cleanup vào, chính hàm cleanup, và một đối số
cho hàm cleanup.
Một khác biệt quan trọng giữa hàm này và `runtime.SetFinalizer` là
hàm cleanup nhận một đối số khác với đối tượng mà chúng ta gắn cleanup vào.
Thay đổi này sửa được một số vấn đề của finalizer.

Không có gì bí mật khi [finalizer
khó dùng đúng](/doc/gc-guide#Common_finalizer_issues).
Ví dụ, các đối tượng được gắn finalizer không được tham gia
vào bất kỳ chu kỳ tham chiếu nào (thậm chí một con trỏ trỏ đến chính nó cũng là quá nhiều), nếu không
đối tượng sẽ không bao giờ được thu hồi và finalizer sẽ không bao giờ chạy, gây ra
rò rỉ.
Finalizer cũng làm trì hoãn đáng kể việc thu hồi bộ nhớ.
Cần tối thiểu hai chu kỳ thu gom rác hoàn chỉnh để thu hồi
bộ nhớ của một đối tượng có finalizer: một lần để xác định rằng nó không còn truy cập được, và
lần tiếp theo để xác định rằng nó vẫn không truy cập được sau khi finalizer
chạy.

Vấn đề là finalizer [hồi sinh đối tượng mà chúng được gắn
vào](https://en.wikipedia.org/wiki/Object_resurrection).
Finalizer không chạy cho tới khi đối tượng không còn truy cập được, lúc đó nó được
xem là "đã chết".
Nhưng vì finalizer được gọi với một con trỏ tới đối tượng, bộ gom rác
phải ngăn việc thu hồi bộ nhớ của đối tượng đó, và thay vào đó phải
tạo ra một tham chiếu mới cho finalizer, khiến nó lại trở nên truy cập được, hay “sống,”
một lần nữa.
Tham chiếu đó thậm chí có thể tồn tại sau khi finalizer trả về, ví dụ nếu
finalizer ghi nó vào một biến toàn cục hoặc gửi nó qua một channel.
Việc hồi sinh đối tượng là vấn đề vì nó có nghĩa là đối tượng, và mọi thứ
nó trỏ tới, và mọi thứ các đối tượng đó trỏ tới, v.v., đều là truy cập được,
ngay cả khi lẽ ra chúng đã bị thu gom như rác.

Chúng tôi giải quyết cả hai vấn đề này bằng cách không truyền đối tượng gốc cho
hàm cleanup.
Trước hết, các giá trị mà đối tượng tham chiếu tới không cần được giữ
ở trạng thái đặc biệt là còn truy cập được bởi bộ gom rác, nên đối tượng vẫn có thể bị thu hồi
ngay cả nếu nó nằm trong một chu kỳ.
Thứ hai, vì đối tượng không cần cho cleanup, bộ nhớ của nó có thể được
thu hồi ngay lập tức.

## Weak pointer

Quay lại ví dụ tệp ánh xạ bộ nhớ của chúng ta, giả sử ta nhận ra rằng chương trình
thường xuyên ánh xạ cùng một tệp lặp đi lặp lại, từ các goroutine khác nhau
không hề biết về nhau.
Điều này ổn từ góc độ bộ nhớ, vì các ánh xạ này sẽ dùng chung
bộ nhớ vật lý, nhưng nó dẫn tới rất nhiều system call không cần thiết để ánh xạ và bỏ ánh xạ tệp.
Điều này đặc biệt tệ nếu mỗi goroutine chỉ đọc một phần nhỏ của mỗi
tệp.

Vậy hãy loại bỏ trùng lặp các ánh xạ theo tên tệp.
(Hãy giả định rằng chương trình của chúng ta chỉ đọc từ các ánh xạ, và các tệp
không bao giờ bị chỉnh sửa hoặc đổi tên sau khi được tạo.
Những giả định như vậy là hợp lý với các tệp phông chữ hệ thống chẳng hạn.)

Chúng ta có thể duy trì một map từ tên tệp tới ánh xạ bộ nhớ, nhưng rồi sẽ không rõ
khi nào an toàn để xóa các phần tử khỏi map đó.
Ta *gần như* có thể dùng cleanup, nếu không vì việc chính bản thân phần tử map
sẽ giữ cho đối tượng tệp ánh xạ bộ nhớ còn sống.

Weak pointer giải quyết vấn đề này.
Weak pointer là một loại con trỏ đặc biệt mà bộ gom rác bỏ qua
khi quyết định một đối tượng còn truy cập được hay không.
[Kiểu weak pointer mới của Go 1.24 là `weak.Pointer`](https://pkg.go.dev/weak#Pointer) có một phương thức `Value`
trả về hoặc là con trỏ thật nếu đối tượng vẫn còn truy cập được, hoặc `nil` nếu không.

Nếu thay vào đó ta duy trì một map chỉ *yếu* trỏ tới tệp ánh xạ bộ nhớ,
ta có thể dọn phần tử map khi không còn ai dùng nó nữa!
Hãy xem điều này trông như thế nào.

```
var cache sync.Map // map[string]weak.Pointer[MemoryMappedFile]

func NewCachedMemoryMappedFile(filename string) (*MemoryMappedFile, error) {
	var newFile *MemoryMappedFile
	for {
		// Try to load an existing value out of the cache.
		value, ok := cache.Load(filename)
		if !ok {
			// No value found. Create a new mapped file if needed.
			if newFile == nil {
				var err error
				newFile, err = NewMemoryMappedFile(filename)
				if err != nil {
					return nil, err
				}
			}

			// Try to install the new mapped file.
			wp := weak.Make(newFile)
			var loaded bool
			value, loaded = cache.LoadOrStore(filename, wp)
			if !loaded {
				runtime.AddCleanup(newFile, func(filename string) {
					// Only delete if the weak pointer is equal. If it's not, someone
					// else already deleted the entry and installed a new mapped file.
					cache.CompareAndDelete(filename, wp)
				}, filename)
				return newFile, nil
			}
			// Someone got to installing the file before us.
			//
			// If it's still there when we check in a moment, we'll discard newFile
			// and it'll get cleaned up by garbage collector.
		}

		// See if our cache entry is valid.
		if mf := value.(weak.Pointer[MemoryMappedFile]).Value(); mf != nil {
			return mf, nil
		}

		// Discovered a nil entry awaiting cleanup. Eagerly delete it.
		cache.CompareAndDelete(filename, value)
	}
}
```

Ví dụ này hơi phức tạp, nhưng ý chính thì đơn giản.
Ta bắt đầu với một map đồng thời toàn cục chứa mọi tệp ánh xạ mà ta đã tạo.
`NewCachedMemoryMappedFile` tra map này để lấy một tệp ánh xạ đã có,
và nếu thất bại thì tạo rồi thử chèn một tệp ánh xạ mới.
Điều này dĩ nhiên cũng có thể thất bại vì ta đang đua với các lần chèn khác, nên
ta cũng phải cẩn thận với điều đó và thử lại.
(Thiết kế này có một khuyết điểm là ta có thể lãng phí bằng cách ánh xạ cùng một tệp nhiều lần
trong một cuộc đua, rồi sẽ phải vứt nó đi thông qua cleanup được thêm bởi
`NewMemoryMappedFile`.
Phần lớn thời gian điều này có lẽ không phải vấn đề lớn.
Cách sửa nó được dành làm bài tập cho người đọc.)

Hãy xem một số thuộc tính hữu ích của weak pointer và cleanup mà đoạn mã này khai thác.

Thứ nhất, hãy lưu ý rằng weak pointer có thể so sánh được.
Không chỉ vậy, weak pointer có danh tính ổn định và độc lập,
vẫn giữ nguyên ngay cả sau khi đối tượng mà chúng trỏ tới đã biến mất từ lâu.
Đó là lý do vì sao hàm cleanup có thể an toàn gọi `CompareAndDelete`
của `sync.Map`, nơi so sánh `weak.Pointer`, và cũng là lý do then chốt khiến
đoạn mã này hoạt động được.

Thứ hai, hãy quan sát rằng ta có thể thêm nhiều cleanup độc lập cho một đối tượng
`MemoryMappedFile`.
Điều này cho phép ta dùng cleanup theo cách có thể kết hợp và dùng chúng để xây
dựng các cấu trúc dữ liệu tổng quát.
Trong ví dụ cụ thể này, có lẽ sẽ hiệu quả hơn nếu kết hợp
`NewCachedMemoryMappedFile` với `NewMemoryMappedFile` để
chúng dùng chung một cleanup.
Tuy nhiên, ưu điểm của đoạn mã ở trên là nó có thể được viết lại
theo cách tổng quát!

```
type Cache[K comparable, V any] struct {
	create func(K) (*V, error)
	m     sync.Map
}

func NewCache[K comparable, V any](create func(K) (*V, error)) *Cache[K, V] {
	return &Cache[K, V]{create: create}
}

func (c *Cache[K, V]) Get(key K) (*V, error) {
	var newValue *V
	for {
		// Try to load an existing value out of the cache.
		value, ok := cache.Load(key)
		if !ok {
			// No value found. Create a new mapped file if needed.
			if newValue == nil {
				var err error
				newValue, err = c.create(key)
				if err != nil {
					return nil, err
				}
			}

			// Try to install the new mapped file.
			wp := weak.Make(newValue)
			var loaded bool
			value, loaded = cache.LoadOrStore(key, wp)
			if !loaded {
				runtime.AddCleanup(newValue, func(key K) {
					// Only delete if the weak pointer is equal. If it's not, someone
					// else already deleted the entry and installed a new mapped file.
					cache.CompareAndDelete(key, wp)
				}, key)
				return newValue, nil
			}
		}

		// See if our cache entry is valid.
		if mf := value.(weak.Pointer[V]).Value(); mf != nil {
			return mf, nil
		}

		// Discovered a nil entry awaiting cleanup. Eagerly delete it.
		cache.CompareAndDelete(key, value)
	}
}
```

## Lưu ý và hướng phát triển tương lai

Bất chấp những nỗ lực tốt nhất của chúng tôi, cleanup và weak pointer vẫn có thể dễ gây lỗi.
Để hướng dẫn những người đang cân nhắc dùng finalizer, cleanup và weak pointer, gần đây chúng tôi
đã cập nhật [hướng dẫn về bộ gom rác](/doc/gc-guide#Finalizers_cleanups_and_weak_pointers) với một số lời khuyên
về cách dùng các tính năng này.
Hãy xem nó vào lần tới khi bạn định dùng chúng, nhưng cũng hãy cân nhắc kỹ liệu
bạn có thật sự cần dùng chúng không.
Đây là những công cụ nâng cao với ngữ nghĩa tinh vi và, như hướng dẫn nói, phần lớn
mã Go hưởng lợi từ các tính năng này một cách gián tiếp chứ không phải từ việc dùng trực tiếp.
Hãy bám vào những trường hợp sử dụng mà chúng thật sự tỏa sáng, và bạn sẽ ổn.

Hiện tại, chúng tôi sẽ nêu ra một số vấn đề mà bạn có khả năng gặp phải hơn.

Thứ nhất, đối tượng mà cleanup được gắn vào không được truy cập được từ
hàm cleanup (dưới dạng biến được capture) cũng như từ đối số của hàm cleanup.
Cả hai tình huống này đều khiến cleanup không bao giờ chạy.
(Trong trường hợp đặc biệt khi đối số cleanup chính là con trỏ được truyền vào
`runtime.AddCleanup`, `runtime.AddCleanup` sẽ panic, như một tín hiệu gửi tới
người gọi rằng họ không nên dùng cleanup theo cách dùng finalizer.)

Thứ hai, khi weak pointer được dùng làm khóa map, đối tượng được tham chiếu yếu
không được phép còn truy cập được từ giá trị map tương ứng, nếu không đối tượng đó
sẽ tiếp tục còn sống.
Điều này có vẻ hiển nhiên khi đang đọc sâu trong một bài viết blog về weak pointer,
nhưng lại là một chi tiết tinh vi rất dễ bỏ lỡ.
Vấn đề này đã truyền cảm hứng cho toàn bộ khái niệm
[ephemeron](https://en.wikipedia.org/wiki/Ephemeron) để giải quyết nó, và đó là
một hướng phát triển tiềm năng trong tương lai.

Thứ ba, một mẫu phổ biến với cleanup là cần một đối tượng bao bọc, giống như
trong ví dụ `MemoryMappedFile`.
Trong trường hợp cụ thể này, bạn có thể hình dung bộ gom rác trực tiếp
theo dõi vùng bộ nhớ được ánh xạ và truyền quanh `[]byte` bên trong.
Chức năng như vậy có thể là một hướng phát triển trong tương lai, và một API cho nó gần đây đã được
[đề xuất](/issue/70224).

Cuối cùng, cả weak pointer lẫn cleanup đều vốn dĩ không xác định, hành vi của chúng
phụ thuộc chặt chẽ vào thiết kế và động lực của bộ gom rác.
Tài liệu cho cleanup thậm chí còn cho phép bộ gom rác không bao giờ chạy
cleanup.
Việc kiểm thử hiệu quả mã dùng chúng có thể khá khó, nhưng [vẫn làm được](/doc/gc-guide#Testing_object_death).

## Tại sao là bây giờ?

Weak pointer đã được nhắc tới như một tính năng cho Go gần như từ những ngày đầu,
nhưng suốt nhiều năm không được đội Go ưu tiên.
Một lý do là chúng tinh vi, và không gian thiết kế của weak pointer là một bãi mìn
gồm các quyết định có thể làm chúng còn khó dùng hơn.
Lý do khác là weak pointer là một công cụ ngách, trong khi đồng thời lại làm tăng
độ phức tạp của ngôn ngữ.
Chúng tôi đã có kinh nghiệm về mức độ khó chịu khi dùng `SetFinalizer`.
Nhưng có một số chương trình hữu ích không thể diễn đạt nếu thiếu chúng, và
package `unique` cùng những lý do tồn tại của nó đã thực sự nhấn mạnh điều đó.

Với generics, bài học rút ra từ finalizer, và những hiểu biết có được từ mọi công trình tuyệt vời
do các nhóm ở ngôn ngữ khác như C# và Java thực hiện kể từ đó, thiết kế cho weak
pointer và cleanup đã nhanh chóng hình thành.
Mong muốn dùng weak pointer với finalizer lại làm nảy sinh thêm câu hỏi,
và vì vậy thiết kế cho `runtime.AddCleanup` cũng nhanh chóng thành hình.

## Lời cảm ơn

Tôi muốn cảm ơn tất cả mọi người trong cộng đồng đã đóng góp phản hồi cho các
issue đề xuất và báo lỗi khi các tính năng này trở nên khả dụng.
Tôi cũng muốn cảm ơn David Chase vì đã cùng tôi suy nghĩ rất kỹ về ngữ nghĩa của weak
pointer, và cảm ơn anh ấy, Russ Cox và Austin
Clements vì đã giúp thiết kế `runtime.AddCleanup`.
Tôi muốn cảm ơn Carlos Amedee vì công việc đưa `runtime.AddCleanup`
vào hiện thực, mài giũa và đưa nó vào Go 1.24.
Và cuối cùng tôi muốn cảm ơn Carlos Amedee và Ian Lance Taylor vì công việc
thay thế `runtime.SetFinalizer` bằng `runtime.AddCleanup` trong toàn bộ
thư viện chuẩn cho Go 1.25.
