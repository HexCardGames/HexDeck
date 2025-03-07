import { writable } from "svelte/store";

interface GameState {
    isLobbyOverlayShown: boolean;
}

const initialState: GameState = {
    isLobbyOverlayShown: false,
};

const gameStore = writable<GameState>(initialState);

export const toggleLobbyOverlay = () => {
    gameStore.update((state) => ({
        ...state,
        isLobbyOverlayShown: !state.isLobbyOverlayShown,
    }));
};

export default gameStore;
