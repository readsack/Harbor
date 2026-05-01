# Harbor
Harbor is a team-management application with support for realtime chatting and kanban boards. It is supposed to be self-hosted.

**NOTE:** IF EVERYTHING SEEMS ZOOMED IN, IT'S BECAUSE I MADE THE UI ON A 2K Res Monitor, and didn't realize that it would look like that until it was too late.

## STACK
1. Backend API: Built With Go, from using the standard HTTP library (no frameworks)
2. Client Application: Build using Tauri & Rust, with SolidJS serving as the frontend framework

## STEPS TO RUN:
1. Go to the server folder and run the main.exe file. This is the backend API, meant to be run on the server computer. (Optionally, change the SECRET value in .env file)
2. Go to the release folder and run the harbor-client.exe file. This is the client side application.
3. Put localhost:8080 in the beginning which is a default url where the server runs.
4. Then Signup, Login to your account.
5. From Here, you can check any invites you've received to join an organization, or create an organization.
6. From the Organization Page, if you've created the organization, then you can invite more people, create teams, and remove people as well.
7. Head to the Teams tab, and create a new team. Inside of the Team page, there is chats (which is the place for the realtime chat channels), kanban (which is the team's kanban board) and the members tab.
8. Create a new chat channel and you can start chatting.
9. Create Colums and Cards in the Kanban Page.

## NOTES:
1. This was incredibly scuffed. I was running low on time before the Flavortown submission, so the frontend UI is not up to the mark.
2. This was my first time using both Go and Tauri.
3. I know it's not that good, but I will polish it in the coming weeks.

