const nameInput = document.getElementById("input_username")
const passInput = document.getElementById("input_password")
const form = document.getElementById("main_form")
const modeChanger = document.getElementById("mode_changer")
const formHeading = document.getElementById("form_heading")


const currHost = (window.location.href)
const pathToUsers = "/api/users/new"
const urlToPostUser = new URL(pathToUsers, currHost) 

const pathToSessions = "/api/session"
const urlToPostSession = new URL(pathToSessions, currHost)

const pathToGetSnonce = pathToSessions + "/snonce"
const urlToSnonce = new URL(pathToGetSnonce, currHost)

const modesEnum = [
    { mode: "login", url: urlToPostSession }, 
    { mode: "registration", url: urlToPostUser}
]
let formMode = modesEnum[0].mode

let urlToPost = urlToPostUser

const cycleLoginModes = (function() {
    let i = 0
    const searchParams = new URLSearchParams(window.location.search)
    for (i; i < modesEnum.length; i++) {
        setFormMode(modesEnum[i])
        modeName = modesEnum[i].mode
        if (searchParams.has(modeName)) {
            break
        }
    }
    i++
    return (function* () {
        while (true) {
            for (i; i < modesEnum.length; i++) {
                setFormMode(modesEnum[i])
                yield
            }
            i = 0
        }
    })()
})()
modeChanger.onclick = () => { cycleLoginModes.next() }

function setFormMode(modeObj) {
    let mName = modeObj.mode
    formMode = mName

    setPostURL(modeObj.url)
    setHeadingString(mName)
}

function setPostURL(url) {
    urlToPost = url
}

function setHeadingString(mode) {
    formHeading.textContent = `user ${mode} form`
}


const MIN_NAME_LEN = 4
const MAX_NAME_LEN = 36
const errNameTooShort = "username must be at least 4 characters long"
const errNameTooLong = "username must be at most 36 characters long"
const errInvalidName = "username shouldn't start with a number"
nameInput.onblur = function() {
    let nameText = nameInput.value

    if (nameText.length < MIN_NAME_LEN) {
        nameInput.setCustomValidity(errNameTooShort)
        return
    }

    if (nameText.length > MAX_NAME_LEN ) {
        nameInput.setCustomValidity(errNameTooLong)
        return
    }

    if (!isValidName(nameText)) {
        nameInput.setCustomValidity(errInvalidName)
        return
    }

    nameInput.setCustomValidity("")
}

const isValidName = function() {
    let regexDoesnStartWithChar = /^\D/
    const controlGroup = [regexDoesnStartWithChar]
    return function(nameStr) {
        return controlGroup.every(regex => regex.test(nameStr))
    }
}()

const MIN_PASSWORD_LENGTH = 8
const errPassTooShort = "password must be at least 8 characters long"
const errPassNotSecure = "password isn't secure (must contain at least 1 digit and 1 char and 1 non-space)"

passInput.onblur = function() {
    let passText = passInput.value
    
    if (passText.length < MIN_PASSWORD_LENGTH) {
        passInput.setCustomValidity(errPassTooShort)
        return 
    }

    if (!isSecureSequence(passText)) {
        passInput.setCustomValidity(errPassNotSecure)
        return
    }

    passInput.setCustomValidity("")
}

const isSecureSequence = (function() {
    let regexHasDigits = /\d/
    let regexHasChars = /\S/
    let regexHasNotOnlyDigits = /\D/
    const controlGroup = [regexHasChars, regexHasDigits, regexHasNotOnlyDigits]
    return function(str) {
        return controlGroup.every(regex => {
            return regex.test(str)
        })
    }
})()

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
        return "operation: success"
    }).then(txtStr => {
        alert(txtStr)
    })

    form.reset()
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
            return "error: authentication conflict"
        default:
            return "error; couldn't register"
    }
}
