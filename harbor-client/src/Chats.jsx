import { createEffect, createSignal, For, onCleanup, onMount, Show } from "solid-js"
import WebSocket from '@tauri-apps/plugin-websocket'
import { fetch } from "@tauri-apps/plugin-http"
import './chat.css'
import { reload } from "@solidjs/router"
function Chats(props) {
    //console.log(props.chats)
    let [trigger, setTrigger] = createSignal(false)
    let [ws, setWS] = createSignal(null)
    let [current, setCurrent] = createSignal(-1)
    let [messages, setMessages] = createSignal([])
    let [mesg, setMesg] = createSignal("")
    let [createChat, setCreateChat] = createSignal(false)
    let [chatName, setChatName] = createSignal("")
    createEffect(async _ => {
        trigger()
        //console.log(current())
        if (current() != -1) {
            //onsole.log(props.chats)

            //console.log(current())
            if (props.url == "") {
                return
            }

            let req = await fetch(new URL("chathistory", props.url), {
                method: "POST",
                headers: {
                    "Cookie": props.JWT,
                    "X-Chat-Key": props.chats[current()].key,
                }
            })
            if (!req.ok) {
                console.log(req)
                return
            }
            else {
                let data = await req.json()
                setMessages(_ => data.messages)
            }
            let url = new URL("msgs", props.url)
            url.protocol = "ws"

            let WS = await WebSocket.connect(url, {
                headers: {
                    "X-Chat-Key": props.chats[current()].key,
                    "Cookie": props.JWT
                }
            })
            console.log(WS)
            setWS(_ => WS)
            const removeListener = ws().addListener((msg) => {
                if (msg.type == "Text") {
                    setMessages(msgs => [...msgs, JSON.parse(msg.data)])
                    console.log(msg)
                }
            });

        }
        document.getElementById("chat-scroll").scrollTop = document.getElementById("chat-scroll").scrollHeight
    })

    onCleanup(async () => {
        if (ws() != null) {
            await ws().disconnect()
        }
    })
    let changeTo = async (index) => {
        if (ws() != null) {
            await ws().disconnect()
            setMessages(_ => [])
            setWS(_ => null)
        }
        setCurrent(_ => index)
    }
    let send = async () => {
        if (ws() == null) {
            return
        }
        await ws().send(JSON.stringify({
            content: mesg(),
        }))
        setMesg(_ => "")

    }
    let createNewChat = async () => {
        let decoder = new TextDecoder()

        if (chatName() == "") {
            return
        }
        let data = new FormData()
        data.append("team_id", props.teamID)
        data.append("name", chatName())
        let req = await fetch(new URL("createchat", props.url), {
            method: "POST",
            headers: {
                "Cookie": props.JWT,

            },
            body: data
        })
        console.log(decoder.decode((await req.body.getReader().read()).value))
        setCreateChat(_ => false)
        trigger(v => !v)
        props.reload()
    }

    return (
        <div className="w-full grid grid-cols-3">
            <Show when={createChat()} >
                <div class="absolute z-10 inset-0 m-auto bg-zinc-900/30 backdrop-blur-lg overlay flex  items-center justify-center">
                    <div className="frm load bg-zinc-900 flex flex-col p-10 shadow-xl">
                        <button class="btn mt-3 w-5 h-5 flex items-center justify-center rounded-lg" onClick={_ => setCreateChat(_ => false)}><i class="bi bi-x"></i></button>
                        <div className="text-3xl font-bold text-center mb-10">Create New Chat</div>
                        <input class="text-2xl font-bold border-b-2 border-zinc-100 shadow-xl" placeholder="Enter New Chat's Name" value={chatName()} onChange={e => setChatName(_ => e.target.value)} />
                        <button class="p-2 text-xl btn m-5 rounded-md" onClick={_ => createNewChat()}>Create</button>

                    </div>
                </div>
            </Show>
            <div className="sideb col-span-1 flex flex-col items-stretch bg-zinc-900 text-zinc-100 border-r-2 border-zinc-300 ">
                <button class="p-2 text-2xl btn m-5 rounded-md" onClick={_ => setCreateChat(_ => true)}>+</button>
                <For each={props.chats}>
                    {(chat, index) =>
                        <div class="w-full p-6 border-zinc-100 border-b-1 load transition-all duration-300" classList={{ "bg-zinc-100": current() == index(), "text-zinc-900": current() == index() }} onClick={async () => { await changeTo(index()) }}>{chat.name}</div>
                    }
                </For>
            </div>
            <div className="col-span-2 flex flex-col load chat-cont h-200">
                <div class="chat-list" id="chat-scroll" >
                    <For each={messages()}>
                        {(msg, index) => (
                            <div className="msg p-5 shadow-xl m-5 load">
                                <div className="font-bold">{msg.username}</div>
                                <div className="text-lg">{msg.content}</div>
                            </div>
                        )}
                    </For>
                </div>
                <div className="flex items-center justify-center p-3 border-t-1 border-zinc-400">
                    <input class="shadow-xl grow p-3 rounded-lg" type="text" value={mesg()} onChange={e => setMesg(_ => e.target.value)} />
                    <button class="btn px-5 py-3 rounded-lg ml-3 bg-zinc-100 text-zinc-900" onClick={send}>Send</button>
                </div>
            </div>
        </div>
    )
}

export default Chats