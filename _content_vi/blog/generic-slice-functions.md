---
title: Các hàm generic vững chắc cho slice
date: 2024-02-22
by:
- Valentin Deleplace
summary: Tránh rò rỉ bộ nhớ trong gói slices.
template: true
---

Gói [slices](/pkg/slices) cung cấp các hàm hoạt động với slice của bất kỳ kiểu nào.
Trong bài viết này, chúng ta sẽ thảo luận cách bạn có thể dùng các hàm đó hiệu quả hơn bằng cách hiểu cách slice được biểu diễn trong bộ nhớ và điều đó ảnh hưởng đến garbage collector ra sao; đồng thời, chúng ta cũng sẽ xem cách gần đây chúng tôi điều chỉnh các hàm này để làm chúng bớt gây ngạc nhiên hơn.

Với [Type parameters](/blog/deconstructing-type-parameters), chúng ta có thể viết các hàm như [slices.Index](/pkg/slices#Index) một lần cho mọi loại slice có phần tử comparable:

```
// Index returns the index of the first occurrence of v in s,
// or -1 if not present.
func Index[S ~[]E, E comparable](s S, v E) int {
	for i := range s {
		if v == s[i] {
			return i
		}
	}
	return -1
}
```

Không còn cần phải hiện thực lại `Index` cho từng kiểu phần tử khác nhau nữa.

Gói [slices](/pkg/slices) chứa rất nhiều helper như vậy để thực hiện các thao tác thường gặp trên slice:

```
	s := []string{"Bat", "Fox", "Owl", "Fox"}
	s2 := slices.Clone(s)
	slices.Sort(s2)
	fmt.Println(s2) // [Bat Fox Fox Owl]
	s2 = slices.Compact(s2)
	fmt.Println(s2)                  // [Bat Fox Owl]
	fmt.Println(slices.Equal(s, s2)) // false
```

Một số hàm mới (`Insert`, `Replace`, `Delete`, v.v.) sẽ sửa đổi slice. Để hiểu chúng hoạt động ra sao, và cách dùng chúng cho đúng, ta cần xem cấu trúc bên dưới của slice.

Một slice là một khung nhìn lên một phần của mảng. [Bên trong](/blog/slices-intro), slice chứa một con trỏ, một độ dài và một dung lượng. Hai slice có thể dùng chung cùng một mảng nền và có thể nhìn những phần chồng lấp lên nhau.

Ví dụ, slice `s` này là một khung nhìn trên 4 phần tử của một mảng có kích thước 6:

{{image "generic-slice-functions/1_sample_slice_4_6.svg" 450}}

Nếu một hàm thay đổi độ dài của slice được truyền vào như tham số, nó cần trả về một slice mới cho người gọi. Mảng nền có thể vẫn giữ nguyên nếu nó không cần tăng kích thước. Điều này giải thích vì sao [append](/blog/slices) và `slices.Compact` trả về một giá trị, còn `slices.Sort`, vốn chỉ sắp xếp lại các phần tử, thì không.

Hãy xét tác vụ xóa một đoạn của slice. Trước thời kỳ generics, cách tiêu chuẩn để xóa đoạn `s[2:5]` khỏi slice `s` là gọi hàm [append](/ref/spec#Appending_and_copying_slices) để chép phần cuối đè lên phần giữa:

```
s = append(s[:2], s[5:]...)
```

Cú pháp này phức tạp và dễ sai, liên quan đến subslice và một tham số variadic. Chúng tôi đã thêm [slices.Delete](/pkg/slices#Delete) để việc xóa phần tử dễ hơn:

```
func Delete[S ~[]E, E any](s S, i, j int) S {
       return append(s[:i], s[j:]...)
}
```

Hàm một dòng `Delete` biểu đạt rõ hơn ý định của lập trình viên. Hãy xét một slice `s` dài 6, capacity 8, chứa các con trỏ:

{{image "generic-slice-functions/2_sample_slice_6_8.svg" 600}}

Lời gọi này xóa các phần tử tại `s[2]`, `s[3]`, `s[4]` khỏi slice `s`:

```
s = slices.Delete(s, 2, 5)
```

{{image "generic-slice-functions/3_delete_s_2_5.svg" 600}}

Khoảng trống ở các chỉ số 2, 3, 4 được lấp bằng cách dời phần tử `s[5]` sang trái, rồi đặt độ dài mới là `3`.

`Delete` không cần cấp phát một mảng mới, vì nó dời phần tử tại chỗ. Giống như `append`, nó trả về một slice mới. Nhiều hàm khác trong gói `slices` cũng theo mẫu này, bao gồm `Compact`, `CompactFunc`, `DeleteFunc`, `Grow`, `Insert` và `Replace`.

Khi gọi các hàm này, ta phải xem slice gốc là không còn hợp lệ, vì mảng nền đã bị sửa đổi. Sẽ là sai lầm nếu gọi hàm mà bỏ qua giá trị trả về:

```
	slices.Delete(s, 2, 5) // incorrect!
	// s still has the same length, but modified contents
```

## Vấn đề về sự sống còn không mong muốn

Trước Go 1.22, `slices.Delete` không sửa các phần tử nằm giữa độ dài mới và độ dài cũ của slice. Dù slice trả về sẽ không bao gồm các phần tử này, “khoảng trống” tạo ra ở cuối slice gốc đã bị mất hiệu lực vẫn tiếp tục giữ chúng. Những phần tử đó có thể chứa các con trỏ tới những đối tượng lớn (một ảnh 20MB), và garbage collector sẽ không giải phóng bộ nhớ gắn với các đối tượng này. Điều đó dẫn tới rò rỉ bộ nhớ và có thể gây ra các vấn đề hiệu năng đáng kể.

Trong ví dụ trên, ta đã xóa thành công các con trỏ `p2`, `p3`, `p4` khỏi `s[2:5]` bằng cách dời một phần tử sang trái. Nhưng `p3` và `p4` vẫn còn hiện diện trong mảng nền, nằm ngoài độ dài mới của `s`. Garbage collector sẽ không thu hồi chúng. Khó thấy hơn là `p5` không phải một trong các phần tử bị xóa, nhưng bộ nhớ của nó vẫn có thể bị rò vì con trỏ `p5` còn nằm ở phần màu xám của mảng.

Điều này có thể gây bối rối cho lập trình viên nếu họ không biết rằng các phần tử “vô hình” vẫn đang dùng bộ nhớ.

Vì thế chúng ta có hai lựa chọn:

* Hoặc giữ nguyên hiện thực hiệu quả của `Delete`. Người dùng sẽ tự đặt các con trỏ lỗi thời về `nil` nếu họ muốn chắc rằng các giá trị được trỏ tới có thể được giải phóng.
* Hoặc thay đổi `Delete` để luôn đặt các phần tử lỗi thời về giá trị zero. Điều này là thêm việc làm, khiến `Delete` hơi kém hiệu quả hơn. Việc zero hóa con trỏ (đặt chúng về `nil`) cho phép garbage collector thu hồi các đối tượng, khi chúng không còn reachable bởi cách nào khác.

Không dễ thấy lựa chọn nào là tốt hơn. Cách thứ nhất cho hiệu năng theo mặc định, còn cách thứ hai cho sự tiết kiệm bộ nhớ theo mặc định.

## Bản sửa

Một quan sát quan trọng là “đặt các con trỏ lỗi thời về `nil`” không dễ như tưởng tượng. Thật ra, tác vụ này dễ sai đến mức chúng ta không nên đặt gánh nặng đó lên người dùng. Vì tính thực dụng, chúng tôi đã chọn sửa hiện thực của năm hàm `Compact`, `CompactFunc`, `Delete`, `DeleteFunc`, `Replace` để “dọn sạch phần đuôi”. Một hệ quả phụ dễ chịu là tải nhận thức giảm xuống và giờ người dùng không còn phải lo lắng về các rò rỉ bộ nhớ kiểu này nữa.

Trong Go 1.22, đây là hình dạng bộ nhớ sau khi gọi Delete:

{{image "generic-slice-functions/4_delete_s_2_5_nil.svg" 600}}

Phần mã thay đổi trong năm hàm đó dùng hàm dựng sẵn mới [clear](/pkg/builtin#clear) (Go 1.21) để đặt các phần tử lỗi thời về zero value của kiểu phần tử của `s`:

{{image "generic-slice-functions/5_Delete_diff.png" 800}}

Zero value của `E` là `nil` khi `E` là kiểu con trỏ, slice, map, chan hoặc interface.

## Kiểm thử bị lỗi

Thay đổi này đã dẫn đến việc một số kiểm thử từng pass trong Go 1.21 nay bị fail trong Go 1.22 khi các hàm của gói slices được dùng sai. Đây là tin tốt. Khi bạn có bug, kiểm thử nên cho bạn biết.

Nếu bạn bỏ qua giá trị trả về của `Delete`:

```
slices.Delete(s, 2, 3)  // !! INCORRECT !!
```

thì bạn có thể lầm tưởng rằng `s` không chứa con trỏ nil nào. [Ví dụ trong Go Playground](/play/p/NDHuO8vINHv).

Nếu bạn bỏ qua giá trị trả về của `Compact`:

```
slices.Sort(s) // correct
slices.Compact(s) // !! INCORRECT !!
```

thì bạn có thể lầm tưởng rằng `s` đã được sắp xếp và compact đúng cách. [Ví dụ](/play/p/eFQIekiwlnu).

Nếu bạn gán giá trị trả về của `Delete` cho một biến khác, nhưng tiếp tục dùng slice gốc:

```
u := slices.Delete(s, 2, 3)  // !! INCORRECT, if you keep using s !!
```

thì bạn có thể lầm tưởng rằng `s` không chứa con trỏ nil nào. [Ví dụ](/play/p/rDxWmJpLOVO).

Nếu bạn vô tình shadow biến slice, rồi tiếp tục dùng slice gốc:

```
s := slices.Delete(s, 2, 3)  // !! INCORRECT, using := instead of = !!
```

thì bạn có thể lầm tưởng rằng `s` không chứa con trỏ nil nào. [Ví dụ](/play/p/KSpVpkX8sOi).


## Kết luận

API của gói `slices` là một cải tiến thực sự so với cú pháp truyền thống trước generics để xóa hoặc chèn phần tử.

Chúng tôi khuyến khích các lập trình viên dùng các hàm mới này, đồng thời tránh những “bẫy” được liệt kê ở trên.

Nhờ các thay đổi gần đây trong phần hiện thực, một lớp rò rỉ bộ nhớ giờ được tự động tránh, không cần thay đổi gì ở API và cũng không cần thêm công sức từ lập trình viên.


## Đọc thêm

Chữ ký của các hàm trong gói `slices` chịu ảnh hưởng mạnh bởi đặc thù của cách slice được biểu diễn trong bộ nhớ. Chúng tôi khuyên bạn nên đọc

*   [Go Slices: usage and internals](/blog/slices-intro)
*   [Arrays, slices: The mechanics of 'append'](/blog/slices)
*   Cấu trúc dữ liệu [dynamic array](https://en.wikipedia.org/wiki/Dynamic_array)
*   [Tài liệu](/pkg/slices) của gói slices

[Đề xuất ban đầu](/issue/63393) về việc zero hóa các phần tử lỗi thời có rất nhiều chi tiết và bình luận.
