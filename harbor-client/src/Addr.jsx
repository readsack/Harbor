import { load } from "@tauri-apps/plugin-store"
import { createSignal, onMount } from "solid-js"
import { fetch } from "@tauri-apps/plugin-http"
import { reload, useNavigate } from "@solidjs/router"

function Addr(props) {
  let [addr, setAddr] = createSignal("http://")
  let [err, setErr] = createSignal("")
  let store
  let nav = useNavigate()
  onMount(async () => {
    store = await load("store.json", { autoSave: false })
  })

  let handleAddr = async (_) => {
    try {
      let url = URL.parse("ping", addr())

      const res = await fetch(url)

      if (!res.ok) {
        throw new Error()
      }
      await store.set("addr", addr())
      await store.save()
      props.onDataLoaded()
    }
    catch (err) {
      console.log(err)
      setErr(e => "Invalid URL Entered")
    }
  }


  return (
    <div class="w-max h-max bg-zinc-900 p-10 rounded-lg load">
      <div className="flex flex-col justify-center items-center">
        <div className="text-2xl text-center m-5">Enter The Server URL</div>
        <input id="input" type="text" className="w-100 p-5 border-2 border-zinc-600 rounded-md" placeholder="Enter Server URL" value={addr()} onChange={(e) => {
          setAddr(_ => e.target.value)
        }} />
        <button className="btn py-3 px-5 border-2 border-zinc-400 rounded-xl mt-5 w-full" onClick={handleAddr}>Save</button>

      </div>
      <div>{err()}</div>
    </div>
  )
}

export default Addr
