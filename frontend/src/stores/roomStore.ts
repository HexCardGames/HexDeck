import { writable } from 'svelte/store';
import { sessionStore } from './sessionStore';

export const loading = writable<"join" | "create" | false>(false);
export const join_error = writable<string | false>(false);
export const create_error = writable<string | false>(false);
export const rejoinRoomCode = writable<string>("");
export const rejoinRoomSessionData = writable<{ sessionToken: string; userId: string }>({ sessionToken: "", userId: "" });

export async function requestJoinRoom(joinCode: string) {
    loading.set("join");
    join_error.set(false);

    try {
        const controller = new AbortController();
        const timeout = setTimeout(() => controller.abort(), 5000);

        const response = await fetch(`/api/room/join`, {
            method: "POST",
            body: JSON.stringify({
                JoinCode: joinCode.replaceAll("-", ""),
                UsernameProposal: "UsernameProposal",
            }),
            headers: {
                "Content-Type": "application/json",
            },
            signal: controller.signal,
        });

        clearTimeout(timeout);

        if (!response.ok) {
            const data: { StatusCode: string; Message: string } = await response.json();
            if (["invalid_join_code"].includes(data?.StatusCode)) {
                join_error.set("no_room_found");
            } else if (data?.Message) {
                join_error.set(data?.Message);
            } else {
                throw new Error("Server error");
            }
            return;
        }

        const data: { SessionToken: string; PlayerId: string; Username: string; Permissions: any } = await response.json();
        const SessionToken = data.SessionToken;
        const UserId = data.PlayerId;
        joinSession(SessionToken, UserId);
    } catch (error: any) {
        if (error.name === "AbortError") {
            join_error.set("timeout");
        } else {
            join_error.set("request_failed");
        }
        console.error("Error joining room: ", error);
    } finally {
        loading.set(false);
    }
}

export async function requestCreateRoom() {
    loading.set("create");
    create_error.set(false);

    try {
        const controller = new AbortController();
        const timeout = setTimeout(() => controller.abort(), 5000);

        const response = await fetch(`/api/room/create`, {
            method: "POST",
            body: JSON.stringify({
                UsernameProposal: "UsernameProposal",
            }),
            headers: {
                "Content-Type": "application/json",
            },
            signal: controller.signal,
        });

        clearTimeout(timeout);

        if (!response.ok) {
            throw new Error("Server error");
        }

        const data: { SessionToken: string; PlayerId: string; Username: string; Permissions: any } = await response.json();
        const SessionToken = data.SessionToken;
        const UserId = data.PlayerId;
        sessionStore.connect(SessionToken, UserId);
    } catch (error: any) {
        if (error.name === "AbortError") {
            create_error.set("timeout");
        } else {
            create_error.set(String(error));
        }
        console.error("Error creating room:", error);
    } finally {
        loading.set(false);
    }
}

export function joinSession(sessionToken: string, userId: string) {
    try {
        sessionStore.connect(sessionToken, userId);
    } catch (error: any) {
        join_error.set("request_failed");
        console.error("Error joining room session: ", error);
    } finally {
        loading.set(false);
    }
}

export async function checkSessionToken(sessionToken: string | undefined): Promise<boolean> {
    if (!sessionToken) return false;
    const params = new URLSearchParams({ sessionToken: sessionToken });
    const res = await fetch(`/api/check/session?${params}`);
    return res.status == 200;
}

export async function checkJoinCode(joinCode: string | undefined): Promise<boolean> {
    if (!joinCode) return false;
    const params = new URLSearchParams({ JoinCode: joinCode });
    const res = await fetch(`/api/check/joinCode?${params}`);
    return res.status == 200;
}

export async function checkSessionData() {
    const currentSessionData: { sessionToken?: string; userId?: string; joinCode?: string } = JSON.parse(localStorage.getItem("currentSessionIds") || "{}");
    if (await checkSessionToken(currentSessionData.sessionToken)) {
        rejoinRoomSessionData.set(currentSessionData as any);
        return;
    }
    if (await checkJoinCode(currentSessionData.joinCode)) {
        rejoinRoomCode.set(currentSessionData.joinCode as string);
        return;
    }

    const lastSessionData: { sessionToken?: string; userId?: string; joinCode?: string } = JSON.parse(localStorage.getItem("lastSessionIds") || "{}");
    if (await checkSessionToken(lastSessionData.sessionToken)) {
        rejoinRoomSessionData.set(lastSessionData as any);
        return;
    }
    if (await checkJoinCode(lastSessionData.joinCode)) {
        rejoinRoomCode.set(lastSessionData.joinCode as string);
        return;
    }
}