import { createSignal, onMount } from "solid-js";
import "./App.css";
import { load } from "@tauri-apps/plugin-store";
import { useNavigate } from "@solidjs/router";
import { fetch } from "@tauri-apps/plugin-http"


function App() {
  let [addr, setAddr] = createSignal("")
  let [jwt, setJWT] = createSignal("")
  let navigate = useNavigate()
  let store
  onMount(async () => {
    store = await load("store.json", {autoSave: false})
    let address = await store.get("addr")
    //console.log(address)
    if(address == "" || address == undefined){
      navigate("/addr")
    }
    else setAddr(_ => address)
    let JWT = await store.get('jwt')
    //console.log(JWT)
    if(JWT == "" || JWT == undefined){
      navigate("/login")
    }
    else setJWT(_ => JWT)
    let url = new URL("api/user", addr())
    let user_req = await fetch(url, {
      method: "POST",
      headers: {
        "Cookie": JWT
      }
    })
    let user_data = await user_req.json()
    if(!user_data.org_id.Valid){
      navigate("/org")
    }
  })  

  return (
    <>
      <button class="btn" onClick={
        () => {
          store.set("jwt", "")
        }
      }
      >Clear</button>
    </>
    //<button className="btn">{addr()}</button>
  )
}

export default App;
