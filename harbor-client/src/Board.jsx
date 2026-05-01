import { fetch } from "@tauri-apps/plugin-http"
import { createEffect, createSignal, Show } from "solid-js"

function Board(props) {
    let [board, setBoard] = createSignal({ columns: [], id: -1, team_id: -1 })
    let [trigger, setTrigger] = createSignal(false)
    let [createCol, setCreateCol] = createSignal(false)
    let [colName, setColName] = createSignal("")
    let [currentCol, setCurrentCol] = createSignal(-1)
    let [cardVal, setCardVal] = createSignal("")
    let [createItem, setCreateItem] = createSignal(false)
    createEffect(async () => {
        trigger()
        let board_url = new URL("api/team/board", props.addr)
        let req = await fetch(board_url, {
            method: "POST",
            headers: {
                "Cookie": props.JWT
            },
            body: JSON.stringify({
                team_id: props.teamID
            })
        })
        if (req.ok) {
            let data = await req.json()
            console.log(data)
            setBoard(_ => data)
        }
    })

    let createColumn = async () => {
        let decoder = new TextDecoder()
        let sendingData = new FormData()
        sendingData.append("isCol", "1")
        sendingData.append("name", colName())
        sendingData.append("team_id", props.teamID)
        let url = new URL("addboard", props.addr)
        let req = await fetch(url, {
            method: "POST",
            headers: {
                "Cookie": props.JWT
            },
            body: sendingData
        })
        console.log(decoder.decode((await req.body.getReader().read()).value))
        setTrigger(v => !v)
        setCreateCol(_ => false)
    }

    let createCard = async (id) => {
        let decoder = new TextDecoder()
        let sendingData = new FormData()
        sendingData.append("isCol", "0")
        sendingData.append("content", cardVal())
        sendingData.append("col_id", id)
        let url = new URL("addboard", props.addr)
        let req = await fetch(url, {
            method: "POST",
            headers: {
                "Cookie": props.JWT
            },
            body: sendingData
        })
        console.log(decoder.decode((await req.body.getReader().read()).value))
        setTrigger(v => !v)
        setCreateItem(_ => false)
    }

    let deleteCol = async (id) => {
        let data = new FormData()
        let decoder = new TextDecoder()

        data.append("isCol", "1")
        data.append("col_id", id)
        let url = new URL("delboard", props.addr)
        let req = await fetch(url, {
            method: "POST",
            headers: {
                "Cookie": props.JWT
            },
            body: data
        })
        setTrigger(v => !v)
        //console.log(decoder.decode((await req.body.getReader().read()).value))

    }
    let deleteCard = async (id) => {
        let data = new FormData()
        let decoder = new TextDecoder()

        data.append("isCol", "0")
        data.append("card_id", id)
        let url = new URL("delboard", props.addr)
        let req = await fetch(url, {
            method: "POST",
            headers: {
                "Cookie": props.JWT
            },
            body: data
        })
        setTrigger(v => !v)
        //console.log(decoder.decode((await req.body.getReader().read()).value))

    }

    return (
        <div className="w-full grid">
            <div class="flex flex-col min-w-0">
                <div className="text-3xl m-5">Kanban Board</div>
                <Show when={board().id != -1} fallback={(
                    <div class="flex w-full grow items-center justify-center">
                        Loading...
                    </div>
                )}>
                    <div className="p-5 flex w-full grow overflow-scroll min-w-0 load">
                        <button class="p-5 text-3xl btn" onClick={_ => setCreateCol(_ => true)}>+</button>
                        <For each={board().columns}>
                            {(col, index) => (

                                <div className="p-8 shadow-xl h-min m-5 flex flex-col min-w-70 load">
                                    <button class="btn w-6 h-6 flex items-center justify-center rounded-lg mb-6" onClick={_ => deleteCol(col.id)}><i class="bi bi-trash-fill"></i></button>
                                    <div className="text-2xl font-bold overflow-hidden">{col.name}</div>
                                    <button className="btn px-6 py-3 font-bold mt-3 rounded-md" onClick={() => {
                                        setCurrentCol(_ => index)
                                        setCreateItem(_ => true)
                                    }}>+</button>
                                    <For each={col.cards}>
                                        {(card, i) => (
                                            <div class="flex flex-col p-5 shadow-xl my-3 load">
                                                <button class="btn w-6 h-6 flex items-center justify-center rounded-lg mb-6" onClick={_ => deleteCard(card.id)}><i class="bi bi-trash-fill"></i></button>

                                                <div className="text-md font-bold">{card.content}</div>
                                                <div className="text-right text-sm mt-2">By: {card.user.username}</div>
                                            </div>
                                        )}
                                    </For>
                                    <Show when={createItem() && currentCol() == index}>
                                        <button class="btn mt-3 w-5 h-5 flex items-center justify-center rounded-lg" onClick={_ => setCreateItem(_ => false)}><i class="bi bi-x"></i></button>
                                        <input type="text" class="text-xl font-bold border-b-2" placeholder="Enter Card Content" value={cardVal()} onChange={e => setCardVal(_ => e.target.value)} />
                                        <button className="btn px-6 py-3 font-bold mt-3" onClick={() => {
                                            createCard(col.id)
                                        }}>Save</button>
                                    </Show>
                                </div>
                            )}
                        </For>
                        <Show when={createCol()}>
                            <div className="p-5 shadow-xl h-min m-5 flex flex-col">
                                <button class="btn m-3 w-5 h-5 flex items-center justify-center rounded-lg" onClick={_ => setCreateCol(_ => false)}><i class="bi bi-x"></i></button>
                                <input type="text" class="text-xl font-bold border-b-2" placeholder="Enter Column Name" value={colName()} onChange={e => setColName(_ => e.target.value)} />
                                <button className="btn px-6 py-3 font-bold mt-3" onClick={createColumn}>Save</button>

                            </div>
                        </Show>
                    </div>
                </Show>
            </div>
        </div>
    )
}

export default Board