const i_f = document.getElementById("f");
const i_t = document.querySelector("textarea");
const f_l = document.querySelector("label[for='f']");
const f_c = document.querySelector("input.na");

let _a;

i_f.addEventListener("change", function () {
    if (this.files.length > 0) {
        _a = i_t.value;

        f_l.textContent = "上传本地文件";
        f_c.checked = true;
        i_t.value = this.files[0].name;
        i_t.readOnly = true;
    } else {
        f_l.textContent = "文本框输入";
        f_c.checked = false;
        i_t.value = _a;
        i_t.readOnly = false;
    }
});

document.body.addEventListener("paste", (e) => {
    const p = e.clipboardData.items[0];
    if (p?.kind === "file" && p.type.startsWith("image/")) {
        const _b = new DataTransfer();
        _b.items.add(p.getAsFile());
        i_f.files = _b.files;

        _a = i_t.value;

        f_l.textContent = "剪切板图片";
        f_c.checked = true;
        i_t.value = "image";
        i_t.readOnly = true;
    }
});

const i_o = document.querySelector("select");
const d_a = document.getElementById("d_a");
const b_a = document.getElementById("a");
const b_b = document.getElementById("b");
const s_l = document.getElementById("link");
const s_t = document.getElementById("token");

const sizeLimit = 104857600;

d_a.addEventListener("click", function () {
    const i_a = i_t.value;
    const i_b = i_f.files[0];

    if (!i_a && !i_b) return;

    const f = f_c.checked ? i_b : new File([i_a], "text.txt");

    if (f.size > sizeLimit) {
        i_t.value = "文件太大，必须小于 " + sizeLimit + " b";
        return;
    }

    const form = new FormData();
    form.append("f", f);

    b_a.classList.toggle("none");
    b_b.classList.remove("none");
    s_l.textContent = "正在加载...";
    s_t.textContent = "正在加载...";

    fetch(window.location.href, {
        method: "POST",
        body: form,
        headers: {
            "type": i_o.value,
        }
    })
        .then(response => response.json())
        .then(msg => {
            s_l.textContent = msg.link;
            s_t.textContent = msg.token;
        });
});

const d_b = document.getElementById("d_b");

d_b.addEventListener("click", function () {
    f_l.textContent = "文本框输入";
    f_c.checked = false;
    i_t.value = _a = "";
    i_t.readOnly = false;
    const _b = new DataTransfer();
    i_f.files = _b.files;

    b_a.classList.remove("none");
    b_b.classList.toggle("none");
});
