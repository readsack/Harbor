import { load } from "@tauri-apps/plugin-store"
import { createSignal, onMount } from "solid-js"
import { fetch } from "@tauri-apps/plugin-http"
import { useNavigate } from "@solidjs/router"

function Addr() {
    let [addr, setAddr] = createSignal("http://")
    let [err, setErr] = createSignal("")
    let store
    let nav = useNavigate()
    onMount(async () => {
        store = await load("store.json", { autoSave: false })
    })

    let handleAddr = async (e) => {
        try {
            let url = URL.parse("ping", addr())
            
            const res = await fetch(url)
            
            if (!res.ok) {
                throw new Error()
            }
            await store.set("addr", addr())
            await store.save()
            nav("/")
        }
        catch(err) {
            console.log(err)
            setErr(e => "Invalid URL Entered")
        }
    }


    return (
        <main className="flex items-center justify-center w-screen h-screen flex-col">
            <div className="flex flex-col justify-center items-center">
                <div className="text-xl text-center">Enter The Server URL</div>
                <input type="text" className="input input-lg" placeholder="Enter Server URL" value={addr()} onChange={(e) => {
                    setAddr(_ => e.target.value)
                }} />
                <button className="btn btn-outline btn-wide m-3" onClick={handleAddr}>Save</button>

            </div>
            <div>{err()}</div>
        </main>

    )
}

export default Addr