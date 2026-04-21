import { createEffect, createSignal, For, Match, Switch } from "solid-js";
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
          console.log(jwt())
          let data = await getOrgData(JWT, addr())
          setTeams(_ => data.teams)
          setUsers(_ => data.users)
          setOrgData(_ => data.org)
          console.log(teams())
          setPage("home")
        }

      }
    }

  })

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
                  <div className="text-3xl ml-10 mt-10">Members List</div>
                  <div className="user-cont px-10 mt-10 load">
                    <div className="items grid-cols-3 grid p-5 rounded-xl text-zinc-400">
                      <div className="text-md">Username</div>
                      <div className="text-md">Email</div>
                      <div className="text-md">Role</div>
                    </div>
                    <For each={users()}>
                      {(user, _) => (
                        <div class="grid grid-cols-3 p-5 rounded-xl shadow-lg/40 shadow-zinc-950">
                          <div className="usr text-lg">{user.username}</div>
                          <div className="email text-lg">{user.email}</div>
                          <div className="status text-lg text-sky-400">{user.id == orgData().ceo_id ? "CEO" : "Member"}{user.id == userData().id ? " (You)" : ""}</div>
                        </div>
                      )}
                    </For>
                  </div>
                </Match>
                <Match when={current() == "teams"}>
                  <div className="text-3xl ml-10 mt-10">Teams List</div>
                  <div className="team-cont pl-10 mt-10 load grid grid-cols-3 gap-10">
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
                            setPage("team")
                          }}>View</div>
                        </div>
                      )}
                    </For>
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
          } />
        </Match>
      </Switch>
    </main>

    //<button className="btn">{addr()}</button>
  )
}

export default App;
