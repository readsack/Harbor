import { useNavigate } from "@solidjs/router"
import { fetch } from "@tauri-apps/plugin-http"
import { load } from "@tauri-apps/plugin-store"
import { createSignal, For, Match, onMount, Switch } from "solid-js"


function Org() {
  let decoder = new TextDecoder("utf-8")
  let [invites, setInvites] = createSignal(null)
  let store
  let [addr, setAddr] = createSignal("")
  let [jwt, setJWT] = createSignal("")
  let [cur, setCur] = createSignal("inv")
  let nav = useNavigate()

  onMount(async () => {
    store = await load("store.json", { autoSave: false })
    let address = await store.get("addr")
    setAddr(_ => address)
    let JWT = await store.get("jwt")
    setJWT(_ => JWT)
    // console.log(new URL("api/user/invites", addr), jwt)
    let req = await fetch(new URL("api/user/invites", addr()), {
      method: "POST",
      headers: {
        "Cookie": jwt(),
      }
    })
    let v = await req.json()
    setInvites(_ => v.invites)
    //console.log(invites())
    console.log(v.invites)
    if (!req.ok) console.log(decoder.decode((await req.body.getReader().read()).value))

  })

  let closeInv = async (accept, key) => {
    let u = new URL("/closeinvite", addr())
    u.searchParams.append("accept", accept)
    u.searchParams.append("invite", key)
    try {
      let req = await fetch(u, {
        method: "POST",
        headers: {
          "Cookie": jwt()
        }
      })
      if (!req.ok) {
        throw new Error(decoder.decode((await req.body.getReader().read()).value))
      }
      else nav("/")
    }
    catch (err) {
      console.log(err)
    }
  }

  return (
    <div class="w-full flex justify-center">
      <Switch>
        <Match when={cur() == "inv"}>
          <div class="h-180 w-1/2 bg-zinc-900 load rounded-lg flex flex-col">
            <div className="text-3xl ml-10 mt-10">Invites</div>
            <div className="text-lg ml-10 text-zinc-300">Here Are the Invites You Have Received</div>
            <a className="text-md text-sky-400 ml-10" onClick={() => { setCur(_ => "cre") }}>Want to Create An Organization?</a>
            <div className="invites flex flex-col grow">
              <For each={invites()} fallback={(
                <div class="flex flex-col w-full grow justify-center items-center">
                  <div className="text-lg text-zinc-400">No Items D:</div>
                </div>
              )}>
                {(invite, index) => (
                  <div className="flex m-10 p-5 border-2 border-zinc-500 rounded-lg items-center">
                    <div className="text-md">
                      {invite.org_name}
                    </div>
                    <div className="obj grow"></div>
                    <div className="btn p-3 mx-3 rounded-md" onClick={() => { closeInv("1", invite.invite_key) }}>Accept</div>
                    <div className="btn p-3 rounded-md" onClick={() => { closeInv("0", invite.invite_key) }}>Deny</div>

                  </div>
                )}
              </For>
            </div>
          </div>
        </Match>
        <Match when={cur() == "cre"}>
          <div class="h-180 w-1/3 bg-zinc-900 load">
            <div className="text-3xl ml-10 mt-10">Create A New Organization</div>
            <a className="text-md text-sky-400 ml-10" onClick={() => { setCur(_ => "inv") }}>Want to Check Invites Instead?</a>
          </div>
        </Match>
      </Switch>

    </div>
  )
}

export default Org
