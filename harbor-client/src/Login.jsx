import { useNavigate } from "@solidjs/router"
import { fetch } from "@tauri-apps/plugin-http"
import { load } from "@tauri-apps/plugin-store"
import { createSignal, Switch, Match, onMount, createEffect } from "solid-js"


function Login() {
    let [email, setEmail] = createSignal("")
    let [pass, setPass] = createSignal("")
    let [username, setUsername] = createSignal("")
    let [addr, setAddr] = createSignal("")
    let [errL, setErrL] = createSignal("")
    let [errS, setErrS] = createSignal("")
    let decoder = new TextDecoder("utf-8")
    let nav = useNavigate()
    let store
    //let [fD, setFD] = createSignal(new FormData())

    onMount(async () => {
        store = await load('store.json', { autoSave: false });
        let val = await store.get('addr')
        setAddr(_ => val)
        let jwt = await store.get('jwt')
        //console.log(jwt)
    })
    

    let logIn = async (_) => {
        let currentFormData = new FormData()
        currentFormData.set("pass", pass())
        currentFormData.set("email", email())
        try{
            let url = new URL("login", addr()) 
            //console.log(url)
            let res = await fetch(url, {
                method: "post",
                body: currentFormData,
            })  
            if(!res.ok){    
                throw new Error(decoder.decode((await res.body.getReader().read()).value))
            }
            else{
                //console.log(res.headers.getSetCookie()[0])
                await store.set("jwt", res.headers.getSetCookie()[0])
                await store.save()
                nav("/")    
            }
        }
        catch(err) {
            setErrL(_ => err.toString())
            
        }
    }

    let signUp = async (_) => {
        let currentFormData = new FormData()
        currentFormData.set("username", username())
        currentFormData.set("pass", pass())
        currentFormData.set("email", email())
        if(username() == "" || email() == "" || pass() == ""){
            setErrS("All Fields Need to Be Filled")
            return
        }
        let regex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/
        if(!regex.test(email())){
            setErrS("Enter A Valid Email Address")
            return
        }
        try{
            let url = new URL("signup", addr()) 
            //console.log(url)
            let res = await fetch(url, {
                method: "post",
                body: currentFormData,
            })  
            if(!res.ok){    
                throw new Error(decoder.decode((await res.body.getReader().read()).value))
            }
            else{
                setErrS("Account Created Succesfully. Now Log In")
            }
            //console.log(decoder.decode((await res.body.getReader().read()).value))
        }
        catch(err) {
            setErrS(_ => err.toString())
        }
    }

    return (
        <main className="flex justify-center items-center w-screen h-screen">
            <div className="flex w-screen items-center justify-center">
                <div className="flex flex-col items-center justify-center w-max grow">
                    <label className="input m-2">
                        Email
                        <input type="email" className="text-sm" placeholder="Enter Your Email..." value={email()} onChange={(e) => {
                            setEmail(_ => e.target.value)
                        }} />
                    </label>
                    <label className="input m-2">
                        Password
                        <input type="password" className="text-sm" placeholder="Enter Your Password..." value={pass()} onChange={(e) => {
                            setPass(_ => e.target.value)
                            
                        }} />
                    </label>
                    <button className="btn btn-wide btn-outline m-4" on:click={logIn}>Login</button>
                    <div className="text-md">{errL()}</div>
                </div>
                <div class="divider divider-horizontal">OR</div>
                <div className="flex flex-col items-center justify-center w-max grow">
                    <label className="input m-2">
                        Username
                        <input type="text" className="text-sm" placeholder="Enter Your Username..." value={username()} onChange={(e) => {
                            setUsername(_ => e.target.value)
                        }} />
                    </label>
                    <label className="input m-2">
                        Email
                        <input type="email" className="text-sm" placeholder="Enter Your Email..." value={email()} onChange={(e) => {
                            setEmail(_ => e.target.value)
                        }} />
                    </label>
                    <label className="input m-2">
                        Password
                        <input type="password" className="text-sm" placeholder="Enter Your Password..." value={pass()} onChange={(e) => {
                            setPass(_ => e.target.value)
                        }} />
                    </label>
                    <button className="btn btn-wide btn-outline m-4" onClick={signUp}>Signup</button>
                    <div className="text-md">{errS()}</div>
                </div>
            </div>
        </main>
    )
}

export default Login