function put(event) {
    event.preventDefault();

    const data_file = document.querySelector('input[type=file]').files[0];
    const data_input = document.querySelector('textarea').value;

    const form_data = new FormData();

    if (data_file) {
        upload_file = data_file
    } else {
        if (data_input.trim() === '') {
            alert("NULL");
            return;
        }

        upload_file = new File([data_input], "f", {
            type: "text/plain"
        });
    }

    document.getElementById('sub').disabled = true;
    document.getElementById('link').textContent = 'waiting...';

    form_data.append('f', upload_file);

    fetch(window.location.href, {
        method: 'POST',
        body: form_data
    })
        .then(response => response.text())
        .then(back_value => {
            document.getElementById('link').textContent = back_value;
        });
}
