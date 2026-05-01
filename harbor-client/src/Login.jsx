import { reload, useNavigate } from "@solidjs/router"
import { fetch } from "@tauri-apps/plugin-http"
import { load } from "@tauri-apps/plugin-store"
import { createSignal, Switch, Match, onMount, createEffect } from "solid-js"
import './login.css'


function Login(props) {
  let [email, setEmail] = createSignal("hii")
  let [pass, setPass] = createSignal("")
  let [username, setUsername] = createSignal("")
  let [addr, setAddr] = createSignal("")
  let [errL, setErrL] = createSignal("")
  let [errS, setErrS] = createSignal("")
  let [isLogin, setLogin] = createSignal(true)
  let decoder = new TextDecoder("utf-8")
  let nav = useNavigate()
  let store
  //let [fD, setFD] = createSignal(new FormData())

  onMount(async () => {
    store = await load('store.json', { autoSave: false });
    let val = await store.get('addr')
    setAddr(_ => val)
    //console.log(jwt)
  })


  let logIn = async (_) => {
    let currentFormData = new FormData()
    currentFormData.set("pass", pass())
    currentFormData.set("email", email())
    try {
      let url = new URL("login", addr())
      //console.log(url)
      let res = await fetch(url, {
        method: "post",
        body: currentFormData,
      })
      if (!res.ok) {
        throw new Error(decoder.decode((await res.body.getReader().read()).value))
      }
      else {
        //console.log(res.headers.getSetCookie()[0])
        await store.set("jwt", res.headers.getSetCookie()[0])
        await store.save()
        props.onDataLoaded()

      }
    }
    catch (err) {
      setErrL(_ => err.toString())

    }
  }

  let signUp = async (_) => {
    let currentFormData = new FormData()
    currentFormData.set("username", username())
    currentFormData.set("pass", pass())
    currentFormData.set("email", email())
    if (username() == "" || email() == "" || pass() == "") {
      setErrS("All Fields Need to Be Filled")
      return
    }
    let regex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/
    if (!regex.test(email())) {
      setErrS("Enter A Valid Email Address")
      return
    }
    try {
      let url = new URL("signup", addr())
      //console.log(url)
      let res = await fetch(url, {
        method: "post",
        body: currentFormData,
      })
      if (!res.ok) {
        throw new Error(decoder.decode((await res.body.getReader().read()).value))
      }
      else {
        setErrS("Account Created Succesfully. Now Log In")
      }
      //console.log(decoder.decode((await res.body.getReader().read()).value))
    }
    catch (err) {
      setErrS(_ => err.toString())
    }
  }

  return (
    <div className="h-180 w-max rounded-lg p-10 shadow-xl shadow-zinc-950 bg-zinc-900">
      <div className="forms flex flex-col">
        <div className="selector flex w-full justify-around text-lg">
          <div class="opt1 grow text-center p-5" classList={{ "selected": isLogin() }} on:click={() => { setLogin(_ => true) }}>Login</div>
          <div class="opt1 grow text-center p-5" classList={{ "selected": !isLogin() }} on:click={() => { setLogin(_ => false) }}>Signup</div>
        </div>
        <div className="loginForm grow flex flex-col items-center" classList={{ "vis": isLogin(), "invis": !isLogin() }}>
          <div className="text-3xl m-10 mt-20 font-semibold">Log In To Your Account</div>
          <div className="cont flex flex-col w-100">
            <div className="text-md font-semibold mb-2">EMAIL</div>
            <input id="input" type="text" class="w-100 p-5 border-2 border-zinc-600 rounded-md" placeholder="Enter Your Email..." value={email()} onChange={(e) => { setEmail(_ => e.target.value) }} />
            <div className="text-md font-semibold mt-8 mb-2">PASSWORD</div>
            <input id="input" type="password" class="w-100 p-5 border-2 border-zinc-600 rounded-md" placeholder="Enter Your Password..." value={pass()} onChange={(e) => { setPass(_ => e.target.value) }} />
            <button className="btn py-3 px-5 border-2 border-zinc-400 rounded-md mt-8 w-min" onClick={logIn}>Submit</button>
            <div className="text-md text-center mt-5">{errL()}</div>
          </div>
        </div>
        <div className="signUpform grow flex flex-col items-center" classList={{ "vis": !isLogin(), "invis": isLogin() }}>
          <div className="text-3xl m-10 mt-10 font-semibold">Create A New Account</div>
          <div className="cont flex flex-col w-100">
            <div className="text-md font-semibold">USERNAME</div>
            <input id="input" type="text" class="w-100 p-5 border-2 border-zinc-600 rounded-md" placeholder="Enter Your Username..." value={username()} onChange={(e) => { setUsername(_ => e.target.value) }} />
            <div className="text-md font-semibold mt-5">EMAIL</div>
            <input id="input" type="text" class="w-100 p-5 border-2 border-zinc-600 rounded-md" placeholder="Enter Your Email..." value={email()} onChange={(e) => { setEmail(_ => e.target.value) }} />
            <div className="text-md font-semibold  mt-5 ">PASSWORD</div>
            <input id="input" type="password" class="w-100 p-5 border-2 border-zinc-600 rounded-md" placeholder="Enter Your Password..." value={pass()} onChange={(e) => { setPass(_ => e.target.value) }} />
            <button className="btn py-3 px-5 border-2 border-zinc-400 rounded-xl mt-5 w-min" onClick={signUp}>Create</button>
            <div className="text-md text-center mt-5">{errS()}</div>

          </div>
        </div>
      </div>
    </div>
  )
}

export default Login
