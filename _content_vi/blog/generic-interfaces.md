---
title: Generic interfaces
date: 2025-07-07
by:
- Axel Wagner
tags:
- type parameters
- generics
- interfaces
summary: Việc thêm type parameter vào interface mạnh hơn bạn tưởng rất nhiều
template: true
---

Có một ý tưởng không hề hiển nhiên cho đến khi bạn nghe về nó lần đầu: vì interface bản thân cũng là kiểu, chúng cũng có thể có type parameter.
Ý tưởng này hóa ra mạnh đến bất ngờ khi biểu đạt các ràng buộc trên hàm và kiểu generic.
Trong bài viết này, chúng ta sẽ minh họa điều đó bằng cách thảo luận việc dùng interface có type parameter trong một vài tình huống thường gặp.

## Một tree set đơn giản

Để làm ví dụ khởi động, giả sử ta cần một phiên bản generic của [cây tìm kiếm nhị phân](https://en.wikipedia.org/wiki/Binary_search_tree).
Các phần tử lưu trong cây như vậy cần có thứ tự, nên type parameter của ta cần một ràng buộc xác định cách so sánh thứ tự.
Một lựa chọn đơn giản là dùng ràng buộc [cmp.Ordered](/pkg/cmp#Ordered), được giới thiệu trong Go 1.21.
Nó giới hạn type parameter vào các kiểu có thứ tự (chuỗi và số) và cho phép các phương thức của kiểu đó dùng các toán tử so sánh dựng sẵn.

{{raw `
    // The zero value of a Tree is a ready-to-use empty tree.
    type Tree[E cmp.Ordered] struct {
        root *node[E]
    }

    func (t *Tree[E]) Insert(element E) {
        t.root = t.root.insert(element)
    }

    type node[E cmp.Ordered] struct {
        value E
        left  *node[E]
        right *node[E]
    }

    func (n *node[E]) insert(element E) *node[E] {
        if n == nil {
            return &node[E]{value: element}
        }
        switch {
        case element < n.value:
            n.left = n.left.insert(element)
        case element > n.value:
            n.right = n.right.insert(element)
        }
        return n
    }
`}}

([playground](/play/p/H7-n33X7P2h))

Tuy nhiên, cách tiếp cận này có nhược điểm là chỉ hoạt động với các kiểu cơ bản mà toán tử `<` được định nghĩa;
bạn không thể chèn các kiểu struct như [time.Time](/pkg/time#Time).

Ta có thể khắc phục bằng cách yêu cầu người dùng cung cấp một hàm so sánh:

{{raw `
    // A FuncTree must be created with NewFuncTree.
    type FuncTree[E any] struct {
        root *funcNode[E]
        cmp  func(E, E) int
    }

    func NewFuncTree[E any](cmp func(E, E) int) *FuncTree[E] {
        return &FuncTree[E]{cmp: cmp}
    }

    func (t *FuncTree[E]) Insert(element E) {
        t.root = t.root.insert(t.cmp, element)
    }
`}}

([playground](/play/p/tiEjuxCHtFF))

Cách này hoạt động, nhưng cũng có nhược điểm.
Ta không còn có thể dùng zero value của kiểu container nữa, vì nó cần một hàm so sánh được khởi tạo tường minh.
Và việc dùng một trường kiểu hàm khiến trình biên dịch khó inline các lời gọi so sánh hơn, điều này có thể tạo ra overhead runtime đáng kể.

Dùng một phương thức trên kiểu phần tử có thể giải quyết các vấn đề đó, vì phương thức gắn trực tiếp với kiểu.
Phương thức không cần được truyền tường minh và trình biên dịch có thể thấy đích của lời gọi để có thể inline nó.
Nhưng làm sao biểu đạt ràng buộc yêu cầu kiểu phần tử phải cung cấp phương thức cần thiết đó?

## Dùng receiver trong constraint

Cách đầu tiên ta có thể thử là định nghĩa một interface thông thường với phương thức `Compare`:

{{raw `
    type Comparer interface {
        Compare(Comparer) int
    }
`}}

Tuy nhiên, ta sớm nhận ra cách này không tốt.
Để hiện thực interface này, tham số của phương thức tự nó phải là `Comparer`.
Điều đó không chỉ có nghĩa phần hiện thực của phương thức phải type-assert tham số về kiểu của chính nó, mà còn đòi hỏi mọi kiểu phải tham chiếu tường minh đến package của ta với tên `Comparer` (nếu không chữ ký phương thức sẽ không đồng nhất).
Điều đó không thật sự trực giao.

Cách tốt hơn là làm cho chính interface `Comparer` trở thành generic:

{{raw `
    type Comparer[T any] interface {
        Compare(T) int
    }
`}}

`Comparer` này giờ mô tả cả một họ interface, một interface cho mỗi kiểu mà `Comparer` được khởi tạo với.
Một kiểu hiện thực `Comparer[T]` đang nói rằng “tôi có thể so sánh chính mình với một `T`”.
Ví dụ, `time.Time` tự nhiên hiện thực `Comparer[time.Time]` vì [nó có một phương thức `Compare` khớp](/pkg/time#Time.Compare):

{{raw `
    // Implements Comparer[Time]
    func (t Time) Compare(u Time) int
`}}

Điều đó tốt hơn, nhưng vẫn chưa đủ.
Điều ta thật sự muốn là một ràng buộc nói rằng type parameter có thể được so sánh với *chính nó*: ta muốn ràng buộc tự tham chiếu.
Điểm tinh tế ở đây là tính tự tham chiếu đó không cần phải nằm trong chính định nghĩa của interface; cụ thể, ràng buộc cho `T` trong kiểu `Comparer` chỉ là `any`.
Thay vào đó, nó là hệ quả của cách ta dùng `Comparer` làm constraint cho type parameter của `MethodTree`:

{{raw `
    // The zero value of a MethodTree is a ready-to-use empty tree.
    type MethodTree[E Comparer[E]] struct {
        root *methodNode[E]
    }
`}}

([playground](/play/p/LuhzYej_2SP))

Vì `time.Time` hiện thực `Comparer[time.Time]`, giờ nó là một type argument hợp lệ cho container này, và ta vẫn có thể dùng zero value như một container rỗng:

{{raw `
    var t MethodTree[time.Time]
    t.Insert(time.Now())
`}}

Để linh hoạt tối đa, một thư viện có thể cung cấp cả ba phiên bản API.
Nếu muốn giảm trùng lặp, mọi phiên bản có thể dùng chung một hiện thực.
Ta có thể dùng phiên bản dựa trên hàm cho điều đó, vì nó là tổng quát nhất:

{{raw `
    // Insert inserts element into the tree, if E implements cmp.Ordered.
    func (t *Tree[E]) Insert(element E) {
        t.root = t.root.insert(cmp.Compare[E], element)
    }

    // Insert inserts element into the tree, using the provided comparison function.
    func (t *FuncTree[E]) Insert(element E) {
        t.root = t.root.insert(t.cmp, element)
    }

    // Insert inserts element into the tree, if E implements Comparer[E].
    func (t *MethodTree[E]) Insert(element E) {
        t.root = t.root.insert(E.Compare, element)
    }
`}}

([playground](/play/p/jzmoaH5eaIv))

Một quan sát quan trọng ở đây là hiện thực dùng chung (biến thể dựa trên hàm) không bị ràng buộc theo cách nào cả.
Nó phải giữ tính linh hoạt tối đa để phục vụ như một lõi chung.
Ta cũng không lưu hàm so sánh trong một trường của struct.
Thay vào đó, ta truyền nó như tham số vì đối số hàm dễ được trình biên dịch phân tích hơn trường struct.

Dĩ nhiên vẫn có một ít boilerplate.
Mọi hiện thực export đều phải sao chép lại toàn bộ API với các kiểu gọi hơi khác nhau.
Nhưng phần này khá thẳng thắn để viết và để đọc.

## Kết hợp phương thức và type set

Ta có thể dùng cấu trúc cây mới để hiện thực một ordered set, cho phép tra cứu phần tử theo thời gian logarithmic.
Giờ hãy tưởng tượng ta cần cho việc tra cứu chạy trong thời gian hằng số; ta có thể thử làm điều đó bằng cách duy trì thêm một map Go thông thường song song với cây:

{{raw `
    type OrderedSet[E Comparer[E]] struct {
        tree     MethodTree[E] // for efficient iteration in order
        elements map[E]bool    // for (near) constant time lookup
    }
`}}

([playground](/play/p/TANUnnSnDqf))

Tuy nhiên, biên dịch đoạn mã này sẽ cho lỗi:

> invalid map key type E (missing comparable constraint)

Thông báo lỗi cho ta biết rằng cần ràng buộc type parameter mạnh hơn để có thể dùng nó làm khóa của map.
Constraint `comparable` là một constraint đặc biệt được khai báo sẵn, thỏa bởi mọi kiểu mà các toán tử so sánh bằng `==` và `!=` được định nghĩa.
Trong Go, đó cũng chính là tập kiểu có thể được dùng làm khóa cho kiểu `map` dựng sẵn.

Ta có ba cách để thêm constraint này vào type parameter, mỗi cách có đánh đổi riêng:

1. Ta có thể [embed](/ref/spec#Embedded_interfaces) `comparable` vào định nghĩa `Comparer` gốc ([playground](/play/p/g8NLjZCq97q)):

   {{raw `
        type Comparer[E any] interface {
            comparable
            Compare(E) int
        }
   `}}

   Nhược điểm là điều đó cũng khiến các kiểu `Tree` chỉ dùng được với những kiểu `comparable`.
   Nói chung, ta không muốn hạn chế kiểu generic nhiều hơn mức cần thiết.
2. Ta có thể thêm một định nghĩa constraint mới ([playground](/play/p/Z2eg4X8xK5Z)).

   {{raw `
        type ComparableComparer[E any] interface {
            comparable
            Comparer[E]
        }
   `}}

   Cách này gọn, nhưng lại đưa thêm một định danh mới (`ComparableComparer`) vào API, mà chuyện đặt tên thì rất khó.
3. Ta có thể thêm constraint trực tiếp vào kiểu bị ràng buộc chặt hơn ([playground](/play/p/ZfggVma_jNc)):

   {{raw `
        type OrderedSet[E interface {
            comparable
            Comparer[E]
        }] struct {
            tree     Tree[E]
            elements map[E]struct{}
        }
   `}}

   Cách này có thể hơi khó đọc, nhất là nếu phải lặp lại thường xuyên.
   Nó cũng khiến việc tái sử dụng constraint ở nơi khác khó hơn.

Dùng cách nào là một lựa chọn về phong cách và cuối cùng phụ thuộc vào sở thích cá nhân.

## Generic interface nên ràng buộc hay không?

Tới đây, đáng để bàn về constraint trên generic interface.
Bạn có thể muốn định nghĩa một interface cho một kiểu container generic.
Ví dụ, giả sử bạn có một thuật toán cần một cấu trúc dữ liệu kiểu set.
Có nhiều kiểu hiện thực set khác nhau với những đánh đổi khác nhau.
Định nghĩa một interface cho những thao tác set mà bạn cần sẽ làm package linh hoạt hơn, để người dùng quyết định những đánh đổi nào phù hợp với ứng dụng cụ thể:

{{raw `
    type Set[E any] interface {
        Insert(E)
        Delete(E)
        Has(E) bool
        All() iter.Seq[E]
    }
`}}

Một câu hỏi tự nhiên ở đây là constraint trên interface này nên là gì.
Nếu có thể, type parameter trên generic interface nên dùng `any` làm constraint, để cho phép mọi kiểu tùy ý.

Từ các phần thảo luận trước, lý do hẳn đã rõ:
Những hiện thực cụ thể khác nhau có thể cần các constraint khác nhau.
Tất cả các kiểu `Tree` ta đã xem ở trên, cũng như kiểu `OrderedSet`, đều có thể hiện thực `Set` cho kiểu phần tử của chúng, dù các kiểu này có constraint khác nhau.

Ý nghĩa của interface là để việc hiện thực được để ngỏ cho người dùng.
Vì ta không thể dự đoán người dùng sẽ muốn đặt loại constraint nào lên phần hiện thực của họ, hãy cố để mọi constraint (mạnh hơn `any`) ở hiện thực cụ thể chứ không phải ở interface.

## Pointer receiver

Hãy thử dùng interface `Set` trong một ví dụ.
Xét một hàm loại bỏ phần tử trùng lặp trong một chuỗi:

{{raw `
    // Unique removes duplicate elements from the input sequence, yielding only
    // the first instance of any element.
    func Unique[E comparable](input iter.Seq[E]) iter.Seq[E] {
        return func(yield func(E) bool) {
            seen := make(map[E]bool)
            for v := range input {
                if seen[v] {
                    continue
                }
                if !yield(v) {
                    return
                }
                seen[v] = true
            }
        }
    }
`}}

([playground](/play/p/hsYoFjkU9kA))

Hàm này dùng `map[E]bool` như một set đơn giản của các phần tử `E`.
Do đó, nó chỉ hoạt động với các kiểu `comparable`, tức các kiểu có định nghĩa toán tử so sánh bằng dựng sẵn.
Nếu muốn tổng quát hóa cho mọi kiểu, ta cần thay nó bằng một generic set:

{{raw `
    // Unique removes duplicate elements from the input sequence, yielding only
    // the first instance of any element.
    func Unique[E any](input iter.Seq[E]) iter.Seq[E] {
        return func(yield func(E) bool) {
            var seen Set[E]
            for v := range input {
                if seen.Has(v) {
                    continue
                }
                if !yield(v) {
                    return
                }
                seen.Insert(v)
            }
        }
    }
`}}

([playground](/play/p/FZYPNf56nnY))

Tuy nhiên, cách này không chạy được.
`Set[E]` là một kiểu interface, và biến `seen` sẽ được khởi tạo bằng `nil`.
Ta cần dùng một hiện thực cụ thể của interface `Set[E]`.
Nhưng như đã thấy trong bài, không có hiện thực set tổng quát nào hoạt động cho mọi kiểu `any`.

Ta buộc phải yêu cầu người dùng cung cấp một hiện thực cụ thể, thông qua một type parameter bổ sung:

{{raw `
    func Unique[E any, S Set[E]](input iter.Seq[E]) iter.Seq[E] { ... }
`}}

([playground](/play/p/kjkGy5cNz8T))

Tuy nhiên, nếu khởi tạo nó bằng hiện thực set của ta, ta lại vấp phải vấn đề khác:

{{raw `
    // OrderedSet[E] does not satisfy Set[E] (method All has pointer receiver)
    Unique[E, OrderedSet[E]](slices.Values(s))
    // panic: invalid memory address or nil pointer dereference
    Unique[E, *OrderedSet[E]](slices.Values(s))
`}}

Vấn đề đầu tiên thể hiện rõ trong thông báo lỗi: ràng buộc kiểu của ta nói rằng type argument cho `S` phải hiện thực `Set[E]`.
Mà vì các phương thức trên `OrderedSet` dùng pointer receiver, type argument đó cũng phải là kiểu con trỏ.

Khi cố làm điều đó, ta vấp phải vấn đề thứ hai.
Nguyên nhân là vì trong phần hiện thực ta khai báo một biến:

{{raw `
    var seen S
`}}

Nếu `S` là `*OrderedSet[E]`, biến sẽ được khởi tạo bằng `nil` như trước.
Gọi `seen.Insert` sẽ panic.

Nếu ta chỉ có kiểu con trỏ, ta không thể có một biến hợp lệ của kiểu giá trị.
Và nếu chỉ có kiểu giá trị, ta lại không thể gọi các phương thức con trỏ trên nó.
Hệ quả là ta cần cả kiểu giá trị *lẫn* kiểu con trỏ.
Vì vậy ta phải giới thiệu thêm type parameter `PS` với một constraint mới `PtrToSet`:

{{raw `
    // PtrToSet is implemented by a pointer type implementing the Set[E] interface.
    type PtrToSet[S, E any] interface {
        *S
        Set[E]
    }
`}}

([playground](/play/p/Kp1jJRVjmYa))

Mẹo ở đây là mối liên hệ giữa hai type parameter trong chữ ký hàm qua type parameter bổ sung trên interface `PtrToSet`.
`S` tự nó không có constraint, nhưng `PS` phải có kiểu `*S` và phải có các phương thức ta cần.
Nên trên thực tế, ta đang ràng buộc `S` phải có một số phương thức nhất định, nhưng những phương thức đó lại dùng pointer receiver.

Dù định nghĩa một hàm với kiểu constraint này cần thêm một type parameter, điều quan trọng là điều đó không làm phức tạp mã sử dụng nó:
miễn là type parameter phụ này nằm ở cuối danh sách type parameter, nó [có thể được suy ra](/blog/type-inference):

{{raw `
    // The third type argument is inferred to be *OrderedSet[int]
    Unique[int, OrderedSet[int]](slices.Values(s))
`}}

Đây là một mẫu tổng quát và rất đáng ghi nhớ: để khi bạn gặp nó trong công việc của người khác, hoặc khi muốn dùng nó trong chính mã của mình.

{{raw `
    func SomeFunction[T any, PT interface{ *T; SomeMethods }]()
`}}

Nếu bạn có hai type parameter, trong đó một cái bị ràng buộc là con trỏ tới cái còn lại, constraint đó sẽ đảm bảo những phương thức liên quan dùng pointer receiver.

## Có nên ràng buộc theo pointer receiver?

Đến đây, có thể bạn đang cảm thấy khá quá tải.
Mọi thứ khá phức tạp và có vẻ không hợp lý khi mong đợi mọi lập trình viên Go đều hiểu chuyện gì đang xảy ra trong chữ ký hàm này.
Ta cũng phải đưa thêm vài cái tên nữa vào API.
Khi mọi người từng cảnh báo việc thêm generics vào Go, đây chính là một trong những điều họ lo ngại.

Vì vậy nếu bạn thấy mình bị mắc vào những vấn đề như thế, đáng để lùi lại một bước.
Nhiều khi ta có thể tránh độ phức tạp này bằng cách nghĩ về bài toán theo hướng khác.
Trong ví dụ này, ta xây một hàm nhận `iter.Seq[E]` và trả về một `iter.Seq[E]` với các phần tử duy nhất.
Nhưng để loại trùng, ta phải thu thập các phần tử duy nhất vào một set.
Và vì điều đó đòi hỏi ta phải cấp phát chỗ cho toàn bộ kết quả, ta thật ra không nhận được nhiều lợi ích khi biểu diễn kết quả như một luồng.

Nếu nghĩ lại bài toán, ta có thể tránh hoàn toàn type parameter phụ bằng cách dùng `Set[E]` như một giá trị interface thông thường:

{{raw `
    // InsertAll adds all unique elements from seq into set.
    func InsertAll[E any](set Set[E], seq iter.Seq[E]) {
        for v := range seq {
            set.Insert(v)
        }
    }
`}}

([playground](/play/p/woZcHodAgaa))

Bằng cách dùng `Set` như một kiểu interface thông thường, việc người gọi phải cung cấp một giá trị hợp lệ của hiện thực cụ thể trở nên rõ ràng.
Đây là một mẫu rất phổ biến.
Và nếu họ cần một `iter.Seq[E]`, họ chỉ việc gọi `All()` trên `set` để lấy ra một cái.

Điều này có làm người gọi phức tạp hơn đôi chút, nhưng cũng có một lợi thế khác so với constraint pointer receiver:
hãy nhớ rằng ta đã bắt đầu với `map[E]bool` như một kiểu set đơn giản.
Rất dễ hiện thực interface `Set[E]` trên nền tảng đó:

{{raw `
    type HashSet[E comparable] map[E]bool

    func (s HashSet[E]) Insert(v E)       { s[v] = true }
    func (s HashSet[E]) Delete(v E)       { delete(s, v) }
    func (s HashSet[E]) Has(v E) bool     { return s[v] }
    func (s HashSet[E]) All() iter.Seq[E] { return maps.Keys(s) }
`}}

([playground](/play/p/KPPpWa7M93d))

Hiện thực này không dùng pointer receiver.
Do đó, dù hoàn toàn hợp lệ, nó sẽ không dùng được với constraint pointer receiver phức tạp ở trên.
Nhưng nó hoạt động rất tốt với phiên bản `InsertAll`.
Cũng như nhiều constraint khác, việc ép mọi phương thức phải dùng pointer receiver thật ra có thể là quá chặt cho nhiều trường hợp thực tế.

## Kết luận

Hy vọng bài viết này đã minh họa một số mẫu và đánh đổi mà type parameter trên interface cho phép.
Nó là một công cụ mạnh, nhưng cũng đi kèm chi phí.
Những điểm rút ra chính là:

1. Dùng generic interface để biểu đạt constraint trên receiver bằng cách dùng chúng theo kiểu tự tham chiếu.
2. Dùng chúng để tạo các quan hệ bị ràng buộc giữa những type parameter khác nhau.
3. Dùng chúng để trừu tượng hóa qua những hiện thực khác nhau với những kiểu constraint khác nhau.
4. Khi bạn thấy mình rơi vào tình huống phải ràng buộc theo pointer receiver, hãy cân nhắc liệu có thể refactor mã để tránh độ phức tạp phụ này không. Xem ["Có nên ràng buộc theo pointer receivers?"](#co-nen-rang-buoc-theo-pointer-receiver).

Như mọi khi, đừng over-engineer mọi thứ: một lời giải kém linh hoạt hơn nhưng đơn giản và dễ đọc hơn cuối cùng có thể vẫn là lựa chọn khôn ngoan hơn.
