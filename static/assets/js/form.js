const nameInput = document.getElementById("input_username")
const passInput = document.getElementById("input_password")
const form = document.getElementById("registration-form")

const currHost = (window.location.href)
const pathToAPI = "/api/users/new"
const urlToPost = new URL(pathToAPI, currHost) 

form.onsubmit = async function(e) {
    e.preventDefault()

    const formData = new FormData(form)
    const formDataObject = {}
    for (const elem of formData.entries()) {
        formDataObject[elem[0]] = elem[1]
    }
    const jsonFormData = JSON.stringify(formDataObject) 
    await fetch(urlToPost, {
        method: "POST",
        headers: {
            'Content-Type': 'application/json',
        },
        body: jsonFormData,
    }).then(function(res) {
        let msg = res.ok ? "registation: success" : getErrorMessage(res.status) 
        alert(msg)
    })

    window.location.reload()
}

function getErrorMessage(status) {
    switch (status) {
        case 400:
            return "error: malformed data input"
        case 409:
            return "error: username is already taken"
        default:
            return "error; couldn't register"
    }
}
