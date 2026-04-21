import { createSignal } from "solid-js"

function Team(props) {
    let [current, setCurrent] = createSignal("chats")
    return (
        <div class="w-5/6 h-200 bg-zinc-900 grid rounded-lg load grid-cols-4">
            <div className="sdb border-r-2 border-zinc-400 flex flex-col">
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
            <div className="col-span-3"></div>
        </div>
    )
}
export default Team