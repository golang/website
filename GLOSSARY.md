# Translation Glossary

Tệp này xác định các bản dịch tiếng Việt ưu tiên và các quyết định văn phong cho
`_content_vi/`.

Cập nhật bảng thuật ngữ này trong quá trình dịch:

- Thêm thuật ngữ khi một cụm từ lặp lại nhiều lần hoặc mang tính chuyên ngành.
- Ưu tiên một bản dịch thống nhất cho mỗi khái niệm, trừ khi ngữ cảnh đòi hỏi khác đi.
- Giữ nguyên tên sản phẩm, tên kho lưu trữ, tên tệp và cú pháp truy vấn.
- Giữ nguyên mã nguồn, URL và các đoạn lệnh.
- Khi chủ ý giữ nguyên một thuật ngữ không dịch, hãy ghi rõ quyết định đó tại đây.

## Quy ước

- Văn phong: rõ ràng, kỹ thuật, trung tính.
- Tên sản phẩm: `Go`
- Tên dự án: `dự án Go`
- Kho lưu trữ: `kho lưu trữ`
- Lịch sử quản lý mã nguồn: `lịch sử quản lý mã nguồn`
- Tác giả: `tác giả`
- Người đóng góp: `người đóng góp`
- Commit: `commit`
- Thay đổi: `thay đổi`
- Tìm kiếm: `tìm kiếm`
- Chuyển hướng: `chuyển hướng`
- Máy chủ: `máy chủ`
- Nguồn tham chiếu chính: `nguồn thông tin chính xác nhất`

## Thuật ngữ đã chốt

| Thuật ngữ gốc | Bản dịch ưu tiên | Ghi chú |
| --- | --- | --- |
| Authors of Go | Tác giả của Go | Page title in `_content_vi/AUTHORS.md`. |
| Go project | dự án Go | Keep `Go` untranslated. |
| AUTHORS | AUTHORS | Filename, do not translate. |
| CONTRIBUTORS | CONTRIBUTORS | Filename, do not translate. |
| author | tác giả | Use for people credited for commits or changes. |
| contributor | người đóng góp | Use when the source explicitly says contributor. |
| commit | commit | Keep the Git term unchanged. |
| repository | kho lưu trữ | Use in general prose. |
| source control history | lịch sử quản lý mã nguồn | Prefer this over shorter paraphrases. |
| authoritative source | nguồn thông tin chính xác nhất | Matches current translation. |
| Go's Gerrit instance | máy chủ Gerrit của Go | Keep `Gerrit` as a product name. |
| search for changes | tìm kiếm các thay đổi | General action phrase. |
| redirect | chuyển hướng | For front matter semantics and prose. |
| ecosystem | hệ sinh thái | Prefer this in Go platform context. |
| community | cộng đồng | General community term. |
| open-source | mã nguồn mở | Use in general prose. |
| anniversary | kỷ niệm | Use for anniversary or milestone context. |
| birthday | sinh nhật | Use for celebratory wording. |
| cloud software | phần mềm đám mây | Prefer this in infra context. |
| garbage collector | bộ gom rác | Preferred runtime term. |
| tooling | hệ thống công cụ | Use for developer-tooling context. |
| release | bản phát hành | Use for product or language releases. |
| vulnerability | lỗ hổng bảo mật | Security context. |
| dependency | dependency | Keep this term in technical prose when clearer. |
| generics | generics | Keep the feature name unchanged for now. |
| production-ready | sẵn sàng cho môi trường production | Preferred phrase in release notes. |
| editor-neutral | trung lập với trình soạn thảo | For tooling compatibility context. |

## Quyết định đang áp dụng

- Dùng tiếng Việt có dấu nhất quán trong toàn bộ `_content_vi/`.
- Ưu tiên cách diễn đạt tự nhiên trong tiếng Việt thay vì bám sát từng từ tiếng Anh.
