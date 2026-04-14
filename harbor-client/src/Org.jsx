import { useNavigate } from "@solidjs/router"
import { fetch } from "@tauri-apps/plugin-http"
import { load } from "@tauri-apps/plugin-store"
import { createSignal, For, onMount } from "solid-js"


function Org() {
  let decoder = new TextDecoder("utf-8")
  let [invites, setInvites] = createSignal([])
  let store
  let [addr, setAddr] = createSignal("")
  let [jwt, setJWT] = createSignal("")
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
    console.log(v.invites[0])
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
    <main >
      <div className="flex w-screen h-screen">
        <ul className="list flex flex-col grow">
          <div className="text-xl font-bold m-5">Invites You Have Received</div>

          <For each={invites()}>
            {(invite, index) => (
              <li className="text-md list-row shadow-md bg-base-200 p-7 ml-5 flex items-center">
                <div className="grow">{invite.org_name}</div>
                <div className="btns">
                  <button className="btn btn-ghost btn-square" onClick={async () => {
                    await closeInv("1", invite.invite_key)
                  }}>
                    <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" fill="currentColor" class="bi bi-check" viewBox="0 0 16 16">
                      <path d="M10.97 4.97a.75.75 0 0 1 1.07 1.05l-3.99 4.99a.75.75 0 0 1-1.08.02L4.324 8.384a.75.75 0 1 1 1.06-1.06l2.094 2.093 3.473-4.425z" />
                    </svg>
                  </button>
                  <button className="btn btn-square btn-ghost ml-5" onClick={async () => {
                    await closeInv("0", invite.invite_key)
                  }}>
                    <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" fill="currentColor" class="bi bi-x" viewBox="0 0 16 16">
                      <path d="M4.646 4.646a.5.5 0 0 1 .708 0L8 7.293l2.646-2.647a.5.5 0 0 1 .708.708L8.707 8l2.647 2.646a.5.5 0 0 1-.708.708L8 8.707l-2.646 2.647a.5.5 0 0 1-.708-.708L7.293 8 4.646 5.354a.5.5 0 0 1 0-.708" />
                    </svg>
                  </button>
                </div>
              </li>

            )}
          </For>
        </ul>
        <div class="divider divider-horizontal">OR</div>
        <div className="flex grow">

        </div>
      </div>
    </main>
  )
}

export default Org
