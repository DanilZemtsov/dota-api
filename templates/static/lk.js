getUserInfo()

async function getUserInfo () {
    let token = localStorage.getItem(`token`)
    const res = await fetch(`http://localhost:8090/user/`,{
        method: "Get",
        headers:{token: `Bearer `+token}
    })
    let result = await res.json()
    console.log(result)
}
