const nameInput = document.getElementById("input_username")
const passInput = document.getElementById("input_password")
const form = document.getElementById("registration-form")

const currHost = (window.location.href)
const pathToPostAPI = "/api/users/new"
const urlToPost = new URL(pathToPostAPI, currHost) 

const pathToGetSnonce = pathToPostAPI + "/snonce"
const urlToSnonce = new URL(pathToGetSnonce, currHost)

const headerAuthReqId = "Authentication-Request-ID"
const headerAuthServerNonce = "Authentication-Server-Nonce"

form.onsubmit = async function(e) {
    e.preventDefault()

    const formData = new FormData(form)

    const username = formData.get("username")
    const password = formData.get("password")

    const hashBufPass = await hashData(password)
    const [reqid, serverNonceByteArr] = await fetch(urlToSnonce, {
        method: "GET"
    }).then(res => {
        const requestid = res.headers.get(headerAuthReqId)
        const snonceHex = res.headers.get(headerAuthServerNonce)
        const snonceArr = decodeHexToByteArr(snonceHex)
        return [requestid, snonceArr]
    })
        
    const [clientNonce, encryptedPassword] = await encryptData(hashBufPass, serverNonceByteArr)
    
    const formDataObject = {
        username: username,
        password: encryptedPassword,
        cnonce: clientNonce,
    }
    console.log("hash", encodeBufferToHex(hashBufPass))
    console.log(new Uint8Array(hashBufPass), serverNonceByteArr)
    console.log(formDataObject)
    const jsonFormData = JSON.stringify(formDataObject) 
    await fetch(urlToPost, {
        method: "POST",
        headers: {
            'Content-Type': 'application/json',
            [headerAuthReqId]: reqid,
        },
        body: jsonFormData,
    }).then((res) => {
        if (res.status === 400) {
            return res.text()
        }
        if (!res.ok) {
            return getErrorMessage(res.status)
        }
        return "registation: success"
    }).then(txtStr => {
        alert(txtStr)
    })

    window.location.reload()
}

async function hashData(str) {
    const data = new TextEncoder().encode(str)
    const hashBuf = await window.crypto.subtle.digest("SHA-256", data)
    return hashBuf
}

async function encryptData(arrBuffer, secretKey) {
    const algo = "AES-GCM"
    const data = new Uint8Array(arrBuffer)
    const ivNonce = window.crypto.getRandomValues(new Uint8Array(12))

    return window.crypto.subtle.importKey(
        "raw", secretKey, 
        {name: algo}, false,
        ["encrypt"],
    ).then((keyObj) => {
        return window.crypto.subtle.encrypt(
            {
                name: algo,
                iv: ivNonce,
            }, 
            keyObj, data
        )
    }).then((encryptedDataBuffer) => {
        return [
            encodeBufferToHex(ivNonce),
            encodeBufferToHex(encryptedDataBuffer)
        ]
    })

}

function decodeHexToByteArr(hexStr) {
    const destLen = ~~(hexStr.length / 2)
    const destArr = new Uint8Array(destLen)
    for (let i = 0; i < destLen; i++) {
        let strI = i * 2
        let currPair = hexStr[strI] + hexStr[strI+1]
        destArr[i] = parseInt(currPair, 16)
    }
    return destArr
}

function encodeBufferToHex(arrayBuf) {
    let arr = new Uint8Array(arrayBuf)
    hexStr = arr.reduce((acc, el) => {
        return acc + ("0" + el.toString(16)).slice(-2)
    }, "")
    return hexStr
}

function getErrorMessage(status) {
    switch (status) {
        case 409:
            return "error: username is already taken"
        default:
            return "error; couldn't register"
    }
}
