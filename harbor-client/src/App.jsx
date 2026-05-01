import { createEffect, createSignal, For, Match, Show, Switch } from "solid-js";
import "./App.css";
import { load } from "@tauri-apps/plugin-store";
import { reload, useNavigate } from "@solidjs/router";
import { fetch } from "@tauri-apps/plugin-http"
import Addr from "./Addr";
import Org from "./Org";
import Login from "./Login"
import Team from "./Team"

async function getOrgData(jwt, addr) {
  let url = new URL("api/org", addr)
  let req = await fetch(url, {
    method: "POST",
    headers: {
      "Cookie": jwt
    }
  })
  let data = await req.json()
  return data
}

function App() {
  let [trigger, setTrigger] = createSignal(false)
  let [current, setCurrent] = createSignal("members")
  let [addr, setAddr] = createSignal("")
  let [jwt, setJWT] = createSignal("")
  let navigate = useNavigate()
  let [page, setPage] = createSignal("")
  let [userData, setUserData] = createSignal(null)
  let store
  let dataLoaded = () => {
    setTrigger(v => !v)
  }

  let [users, setUsers] = createSignal(null)
  let [teams, setTeams] = createSignal(null)
  let [orgData, setOrgData] = createSignal(null)
  let [teamID, setTeamID] = createSignal(-1)
  let [mem, showMem] = createSignal(false)
  let [makeTeam, showMakeTeam] = createSignal(false)
  let [memberEmail, setMemberEmail] = createSignal("")
  let [isAdmin, setAdmin] = createSignal(false)
  let [teamName, setTeamName] = createSignal("")
  createEffect(async () => {
    trigger()
    store = await load("store.json", { autoSave: false })
    let address = await store.get("addr")
    //console.log(address)
    if (address == "" || address == undefined) {
      setPage(_ => "addr")
    }
    else setAddr(_ => address)
    console.log(addr())
    let JWT = await store.get('jwt')
    //console.log(JWT)
    if (addr() != "") {
      if (JWT == "" || JWT == undefined) {
        console.log("login")
        setPage(_ => "login")
      }
      else setJWT(_ => JWT)
      if (jwt() != "") {
        let url = new URL("api/user", addr())
        let user_req = await fetch(url, {
          method: "POST",
          headers: {
            "Cookie": JWT
          }
        })
        let user_data = await user_req.json()
        setUserData(_ => user_data)
        if (!user_data.org_id.Valid) {
          //navigate("/org")
          setPage("org")
        }
        else {
          //console.log(jwt())
          let data = await getOrgData(JWT, addr())
          setTeams(_ => data.teams)
          setUsers(_ => data.users)
          setOrgData(_ => data.org)
          setAdmin(_ => userData().id == data.org.ceo_id)
          //console.log(teams())
          setPage("home")
        }

      }
    }

  })

  let addMember = async () => {
    let decoder = new TextDecoder()
    let data = new FormData()
    data.append("email", memberEmail())
    let req = await fetch(new URL("sendinvite", addr()), {
      method: "POST",
      headers: {
        "Cookie": jwt()
      },
      body: data
    })
    if (!req.ok) console.log(decoder.decode((await req.body.getReader().read()).value))
    showMem(_ => false)
    setTrigger(v => !v)
  }

  let createTeam = async () => {
    let decoder = new TextDecoder()
    let data = new FormData()
    data.append("name", teamName())
    let req = await fetch(new URL("createteam", addr()), {
      method: "POST",
      headers: {
        "Cookie": jwt()
      },
      body: data
    })
    if (!req.ok) console.log(decoder.decode((await req.body.getReader().read()).value))
    showMakeTeam(_ => false)
    setTrigger(v => !v)
  }

  let deleteAccount = async () => {
    let req = await fetch(new URL("removeacc", addr()), {
      method: "POST",
      headers: {
        "Cookie": jwt()
      },
    })
    if (req.ok) await logOut()
  }

  let deleteOrg = async () => {
    let decoder = new TextDecoder()
    let req = await fetch(new URL("deleteorg", addr()), {
      method: "POST",
      headers: {
        "Cookie": jwt()
      },
      body: JSON.stringify({
        org_id: orgData().id
      })
    })
    if (!req.ok) console.log(decoder.decode((await req.body.getReader().read()).value))
    setTrigger(v => !v)
  }

  let removeUser = async (user_id) => {
    let decoder = new TextDecoder()
    let req = await fetch(new URL("deleteuser", addr()), {
      method: "POST",
      headers: {
        "Cookie": jwt()
      },
      body: JSON.stringify({
        user_id: user_id
      })
    })
    if (!req.ok) console.log(decoder.decode((await req.body.getReader().read()).value))
    setTrigger(v => !v)
  }

  let logOut = async () => {
    let store = await load("store.json", { autoSave: false })
    store.set("jwt", "")
    setTrigger(v => !v)
  }

  return (
    <main className="flex justify-center items-center w-screen h-screen">

      <button type="button" onClick={() => { store.clear() }} class="absolute top-10 left-10 btn py-3 px-5 border-2 border-zinc-400 rounded-md w-min">Clear</button>
      <div className="backg"></div>
      <div className="attbr">Photo by Spenser Sembrat on Unsplash</div>

      <Switch>

        <Match when={page() == "addr"}>
          <Addr onDataLoaded={dataLoaded} />
        </Match>
        <Match when={page() == "login"}>
          <Login onDataLoaded={dataLoaded} />
        </Match>
        <Match when={page() == "org"}>
          <Org onDataLoaded={dataLoaded} />
        </Match>
        <Match when={page() == "home"}>
          <div className="w-5/6 h-200 bg-zinc-900 grid rounded-lg load grid-cols-4">
            <div className="sdb border-r-2 border-zinc-400 h-200">
              <div className="tabs flex flex-col p-5 h-200">
                <button className="btn rounded-md m-3 w-10 h-10" onClick={e => dataLoaded()}><i class="bi bi-arrow-repeat"></i></button>

                <div className="text-lg text-center font-bold p-5 hover:bg-zinc-800 transition-all duration-300 rounded-xl flex items-center" classList={{ sel: current() == "members" }} onClick={_ => setCurrent(_ => "members")}>
                  <i class="bi bi-people font-bold mr-2"></i>
                  Members
                </div>
                <div className="text-lg text-center font-bold p-5 hover:bg-zinc-800 transition-all duration-300 rounded-xl flex items-center" classList={{ sel: current() == "teams" }} onClick={_ => setCurrent(_ => "teams")}>
                  <i class="bi bi-collection font-bold mr-2"></i>
                  Teams
                </div>
                <div class="grow"></div>
                <div className="text-lg text-center font-bold p-5 hover:bg-zinc-800 transition-all duration-300 rounded-xl flex items-center" classList={{ sel: current() == "account" }} onClick={_ => setCurrent(_ => "account")}>
                  <i class="bi bi-person-circle font-bold mr-2"></i>
                  Account
                </div>

              </div>
            </div>
            <div className="col-span-3">
              <Switch>
                <Match when={current() == "members"}>
                  <Show when={mem()} >
                    <div className="text-3xl backdrop-blur-xl z-10 overlay absolute inset-0 m-auto flex items-center justify-center bg-zinc-900/40">
                      <div className="frm load bg-zinc-900 flex flex-col p-10 shadow-xl">
                        <button class="btn mt-3 w-7 h-7 flex items-center justify-center rounded-xl" onClick={_ => showMem(_ => false)}><i class="bi bi-x"></i></button>
                        <div className="text-3xl font-bold text-center mb-10">Invite New Person</div>
                        <input class="text-2xl font-bold border-b-2 border-zinc-100 shadow-xl min-w-150" placeholder="Email..." value={memberEmail()} onChange={e => setMemberEmail(_ => e.target.value)} />
                        <button class="p-2 text-xl btn m-5 rounded-md" onClick={_ => addMember()}>Send</button>
                      </div>
                    </div>
                  </Show>
                  <div className="text-3xl ml-10 mt-10">Members List
                    <Show when={isAdmin()}>
                      <button className="btn p-2 m-3 rounded-md text-3xl" onClick={_ => showMem(_ => true)}>
                        +
                      </button>
                    </Show>
                  </div>
                  <div className="user-cont px-10 mt-10 load">
                    <div className="items grid p-5 rounded-xl text-zinc-400" style={{ "grid-template-columns": "1fr 1fr 1fr 50px" }}>
                      <div className="text-md">Username</div>
                      <div className="text-md">Email</div>
                      <div className="text-md">Role</div>
                    </div>
                    <For each={users()}>
                      {(user, _) => (
                        <div class="grid p-5 rounded-xl shadow-lg/40 shadow-zinc-950" style={{ "grid-template-columns": "1fr 1fr 1fr 50px" }}>
                          <div className="usr text-lg">{user.username}</div>
                          <div className="email text-lg">{user.email}</div>
                          <div className="status text-lg text-sky-400">{user.id == orgData().ceo_id ? "CEO" : "Member"}{user.id == userData().id ? " (You)" : ""}</div>
                          <Show when={user.id != orgData().ceo_id && isAdmin()}>
                            <button className="btn font-bold text-2xl" onClick={_ => removeUser(user.id)}><i class="bi bi-x" /></button>
                          </Show>
                        </div>
                      )}
                    </For>
                  </div>
                </Match>
                <Match when={current() == "teams"}>
                  <Show when={makeTeam()} >
                    <div className="text-3xl backdrop-blur-xl z-10 overlay absolute inset-0 m-auto flex items-center justify-center bg-zinc-900/40">
                      <div className="frm load bg-zinc-900 flex flex-col p-10 shadow-xl">
                        <button class="btn mt-3 w-7 h-7 flex items-center justify-center rounded-xl" onClick={_ => showMakeTeam(_ => false)}><i class="bi bi-x"></i></button>
                        <div className="text-3xl font-bold text-center mb-10">Create New Team</div>
                        <input class="text-2xl font-bold border-b-2 border-zinc-100 shadow-xl min-w-100" placeholder="Enter Team Name" value={teamName()} onChange={e => setTeamName(_ => e.target.value)} />
                        <button class="p-2 text-xl btn m-5 rounded-md" onClick={_ => createTeam()}>Send</button>
                      </div>
                    </div>
                  </Show>
                  <div className="text-3xl ml-10 mt-10">Teams List
                    <Show when={isAdmin()}>
                      <button className="btn p-2 m-3 rounded-md text-3xl" onClick={_ => showMakeTeam(_ => true)}>
                        +
                      </button>
                    </Show>
                  </div>
                  <div className="team-cont pl-10 mt-10 load grid grid-cols-3 gap-10 pr-10">
                    <For each={teams()}>
                      {(team, _) => (
                        <div class="flex flex-col p-8 bg-zinc-900 shadow-lg/40 shadow-zinc-950 max-w-80 min-h-70 transition-all duration-300 rounded-xl hover:shadow-xl/80 hover:shadow-zinc-950">

                          <div className="ic p-5 bg-sky-400 rounded-xl h-max w-max">
                            <i class="bi bi-people-fill text-xl w-max h-max flex"></i>
                          </div>
                          <div className="text-3xl mt-4">{team.name}</div>
                          <div className="grow"></div>
                          <div className="btn px-5 p-2 w-min rounded-lg" onClick={() => {
                            setTeamID(_ => team.id)
                            //console.log(team.id)
                            setPage("team")
                          }}>View</div>
                        </div>
                      )}
                    </For>
                  </div>
                </Match>
                <Match when={current() == "account"}>
                  <div class="p-10 load flex flex-col">
                    <div className="text-3xl font-bold">Account Overview</div>
                    <div className="text-2xl font-bold my-3">Username: {userData().username}</div>
                    <div className="text-2xl font-bold my-3">Email: {userData().email}</div>
                    <div className="text-2xl font-bold my-3">Organization: {orgData().name}</div>
                    <button className="btn text-2xl my-4 font-bold p-3 w-max" onClick={deleteAccount}>Delete Account</button>
                    <button className="btn text-2xl my-4 font-bold p-3 w-max" onClick={logOut}>Log Out</button>
                    <Show when={isAdmin()}>
                      <button className="btn text-2xl my-4 font-bold p-3 w-max" onClick={deleteOrg}>Delete Organization</button>
                    </Show>
                  </div>
                </Match>
              </Switch>
            </div>
          </div>
        </Match>
        <Match when={page() == "team"}>
          <Team navHome={() => {
            setPage(_ => "home")
            dataLoaded()
          }
          } teamID={teamID()} members={users()} isAdmin={isAdmin()} reload={_ => setTrigger(v => !v)} />
        </Match>
      </Switch>
    </main>

    //<button className="btn">{addr()}</button>
  )
}

export default App;
