// localStorage.setItem("jwt-token", "xxzcsadasdasdasd");
// async function getPas (){
//  const res = await fetch(`http://localhost:8090/page/registr`)
//  console.log(res);
// }

async function register(e) {
 var formData = new FormData(e.target);
 const res = await fetch(`http://localhost:8090/registr/`,{
  method: "Post",
  body : formData,
 })
 let result = await res.json()
 localStorage.setItem(`token`,result.token)
 console.log(result.token)
 window.location.replace("http://localhost:8090/page/lk/")
 return false;
}

async function login(e) {
 var formData = new FormData(e.target);
 const res = await fetch(`http://localhost:8090/login/`,{
  method: "Post",
  body : formData,
 })

 let result = await res.json()
 console.log(result)
 localStorage.setItem(`token`,result.token)
 window.location.replace("http://localhost:8090/page/lk/")
 return false;
}

const formRegister = document.querySelector('.form-register');
formRegister?.addEventListener('submit',async (e) => {
 e.preventDefault();
 await register(e)
});

const formLogin = document.querySelector('.form-login');
formLogin?.addEventListener('submit', async (e) => {
 e.preventDefault();

 await login(e)

})
