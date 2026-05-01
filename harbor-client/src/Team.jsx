import { createEffect, createSignal, For, Match, Switch } from "solid-js"
import { load } from "@tauri-apps/plugin-store"
import { fetch } from "@tauri-apps/plugin-http"
import Chats from "./Chats"
import Board from "./Board"

function Team(props) {
    let [current, setCurrent] = createSignal("chats")
    let [trigger, setTrigger] = createSignal(false)
    let [addr, setAddr] = createSignal("")
    let [jwt, setJWT] = createSignal("")
    let [members, setMembers] = createSignal([])
    let [chats, setChats] = createSignal([])
    let [data, setData] = createSignal(null)
    let [showMem, setShowMem] = createSignal(false)
    let [otherMem, setOtherMem] = createSignal([])
    let [newMem, setNewMem] = createSignal("")
    createEffect(async () => {
        trigger()
        let store = await load("store.json", { autoSave: false })
        let address = await store.get("addr")
        setAddr(_ => address)
        let JWT = await store.get("jwt")
        setJWT(_ => JWT)
        let res = await fetch(new URL("api/team", addr()), {
            method: "POST",
            headers: {
                "Cookie": jwt()
            },
            body: JSON.stringify({
                "team_id": props.teamID
            })
        })
        //console.log(props.teamID)
        //let decoder = new TextDecoder('utf-8')
        //console.log(decoder.decode((await res.body.getReader().read()).value))
        if (res.ok) {
            let data = await res.json()
            console.log(data)
            setChats(_ => data.chats)
            setMembers(_ => data.members)
            setData(_ => data)
            let existingMemberIDs = new Set(data.members.map(user => user.id))
            setOtherMem(_ => props.members.filter(mem => !existingMemberIDs.has(mem.id)))
            console.log(otherMem())
        }
    })

    let addMember = async () => {
        let user_id = parseInt(newMem(), 10)

        let decoder = new TextDecoder('utf-8')
        let url = new URL("teamadd", addr())
        let req = await fetch(url, {
            method: "POST",
            headers: {
                "Cookie": jwt()
            },
            body: JSON.stringify({
                "user_id": user_id,
                "team_id": props.teamID
            })
        })
        if (!req.ok) console.log(decoder.decode((await req.body.getReader().read()).value))
        props.reload()
    }

    let removeUser = async (user_id) => {
        let decoder = new TextDecoder('utf-8')
        let url = new URL("removeuser", addr())
        let req = await fetch(url, {
            method: "POST",
            headers: {
                "Cookie": jwt()
            },
            body: JSON.stringify({
                "user_id": user_id,
                "team_id": props.teamID
            })
        })
        if (!req.ok) console.log(decoder.decode((await req.body.getReader().read()).value))
        props.reload()
    }


    return (
        <div class="w-5/6 h-200 bg-zinc-900 grid rounded-lg load grid-cols-4">
            <div className="sdb border-r-2 border-zinc-400 flex flex-col">
                <button className="btn rounded-md m-3 w-10 h-10" onClick={e => props.reload()}><i class="bi bi-arrow-repeat"></i></button>

                <div className="text-md p-3 font-bold hover:bg-zinc-800 w-min m-2 flex rounded-lg" onClick={props.navHome}>
                    <i class="bi bi-arrow-left font-bold mr-2"></i>
                    Home
                </div>
                <div className="tabs flex flex-col p-5 grow">

                    <div className="text-lg text-center font-bold p-5 hover:bg-zinc-800 transition-all duration-300 rounded-xl flex items-center" classList={{ sel: current() == "chats" }} onClick={_ => setCurrent(_ => "chats")}>
                        <i class="bi bi-chat-left-text font-bold mr-2"></i>
                        Chats
                    </div>
                    <div className="text-lg text-center font-bold p-5 hover:bg-zinc-800 transition-all duration-300 rounded-xl flex items-center" classList={{ sel: current() == "kanban" }} onClick={_ => setCurrent(_ => "kanban")}>
                        <i class="bi bi-kanban font-bold mr-2"></i>
                        Kanban
                    </div>
                    <div class="grow"></div>
                    <div className="text-lg text-center font-bold p-5 hover:bg-zinc-800 transition-all duration-300 rounded-xl flex items-center" classList={{ sel: current() == "members" }} onClick={_ => setCurrent(_ => "members")}>
                        <i class="bi bi-people font-bold mr-2"></i>
                        Members
                    </div>
                </div>
            </div>
            <div className="col-span-3 grid grid-cols-1">
                <Switch>
                    <Match when={current() == "chats"}>
                        <Chats chats={chats()} url={addr()} JWT={jwt()} teamID={props.teamID} reload={() => { setTrigger(v => !v) }} />
                    </Match>
                    <Match when={current() == "members"}>
                        <div>
                            <Show when={showMem()} >
                                <div class="absolute z-10 inset-0 m-auto bg-zinc-900/30 backdrop-blur-lg overlay flex items-center justify-center">
                                    <div className="frm load bg-zinc-900 flex flex-col p-10 shadow-xl">
                                        <button class="btn mt-3 w-7 h-7 flex items-center justify-center rounded-xl" onClick={_ => setShowMem(_ => false)}><i class="bi bi-x"></i></button>
                                        <div className="text-3xl font-bold text-center mb-10">Add A New Member</div>
                                        <select class="text-2xl text-zinc-900 font-bold border-b-2 border-zinc-900 shadow-xl" placeholder="Enter New Member's Name" value={newMem()} onChange={e => setNewMem(_ => e.target.value)}>
                                            <For each={otherMem()}>
                                                {(mem, index) => (
                                                    <option value={mem.id}>{mem.username}</option>
                                                )}
                                            </For>
                                        </select>
                                        <button class="p-2 text-xl btn m-5 rounded-md" onClick={_ => addMember()}>Add</button>

                                    </div>
                                </div>
                            </Show>
                            <div className="text-3xl ml-10 mt-10">Members List
                                <Show when={props.isAdmin}>
                                    <button className="btn p-2 m-3 rounded-md text-3xl" onClick={_ => setShowMem(_ => true)}>
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
                                <For each={members()}>
                                    {(user, _) => (
                                        <div class="grid p-5 rounded-xl shadow-lg/40 shadow-zinc-950" style={{ "grid-template-columns": "1fr 1fr 1fr 50px" }}>
                                            <div className="usr text-lg">{user.username}</div>
                                            <div className="email text-lg">{user.email}</div>
                                            <div className="status text-lg text-sky-400">{user.id == data().team.sup_id ? "Supervisor" : "Member"}</div>
                                            <Show when={user.id != data().team.sup_id && props.isAdmin}>
                                                <button className="btn font-bold text-2xl" onClick={_ => removeUser(user.id)}><i class="bi bi-x" /></button>
                                            </Show>
                                        </div>
                                    )}
                                </For>
                            </div>
                        </div>
                    </Match>
                    <Match when={current() == "kanban"}>
                        <Board JWT={jwt()} addr={addr()} teamID={props.teamID} />
                    </Match>
                </Switch>
            </div>
        </div>
    )
}
export default Team