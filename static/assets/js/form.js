const nameInput = document.getElementById("input_username")
const passInput = document.getElementById("input_password")
const form = document.getElementById("registration-form")

const currHost = (window.location.href)
const pathToAPI = "/api/users/new"
const urlToPost = new URL(pathToAPI, currHost) 

form.onsubmit = async function(e) {
    e.preventDefault()

    const formData = new FormData(form)

    const username = formData.get("username")
    const password = formData.get("password")
    const hashBufPass = await hashData(password)
    // TODO: instead of generating random key recieve it from server
    const serverNonce = window.crypto.getRandomValues(new Uint8Array(16))
    const encryptedPassword = await encryptData(hashBufPass, serverNonce)

    const formDataObject = {
        username: username,
        password: encryptedPassword,
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

async function hashData(str) {
    const data = new TextEncoder().encode(str)
    const hashBuf = await window.crypto.subtle.digest("SHA-256", data)
    return hashBuf
}

async function encryptData(str, key) {
    const algo = "AES-GCM"
    const data = new TextEncoder().encode(str)
    let iv = window.crypto.getRandomValues(new Uint8Array(12))

    return window.crypto.subtle.importKey(
        "raw", key, 
        {name: algo}, false,
        ["encrypt"],
    ).then((keyObj) => {
        return window.crypto.subtle.encrypt(
            {
                name: algo,
                iv: iv,
            }, 
            keyObj, data
        )
    }).then(encodeBufferToBase64)


}

function encodeBufferToBase64(arrayBuf) {
    let decoder = new TextDecoder()
    let uriEncoded = encodeURIComponent(decoder.decode(arrayBuf))
    let b64 = window.btoa(uriEncoded)
    return b64
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
