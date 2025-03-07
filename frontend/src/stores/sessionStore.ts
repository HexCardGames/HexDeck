import { writable, get } from 'svelte/store';
import { io, Socket } from 'socket.io-client';

export enum GameState {
  Undefined = -1,
  Lobby,
  Running,
  Ended
}

interface PlayerPermissionObj {
  isHost: boolean;
}

interface GameOptions { }

interface PlayerObj {
  PlayerId: string;
  Username: string;
  Permissions: number;
  IsConnected: boolean;
}

interface SessionData {
  roomId: string | null;
  joinCode: string | null;
  gameOptions: GameOptions;
  players: Array<PlayerObj>;
  cardDeckId: string | null;
  gameState: GameState;
  socket: Socket | null;
  connected: boolean;
  userId: string | null;
  messages: string[];
  sessionToken: string | null;
}

interface RoomInfoObj {
  RoomId: string;
  JoinCode: string;
  TopCard: any;
  GameState: GameState;
  CardDeckId: number;
  Winner?: string;
  Players: PlayerObj[];
}

interface StatusInfoObj {
  IsError: boolean;
  StatusCode: string;
  Message: string;
}

class SessionManager {
  private store = writable<SessionData>({
    roomId: null,
    joinCode: null,
    gameState: -1,
    gameOptions: {},
    players: [],
    cardDeckId: null,
    socket: null,
    connected: false,
    userId: null,
    messages: [],
    sessionToken: null,
  });

  private socket: Socket | null = null;

  constructor() {
    const storedSessionIds = this.getStoredSessionIds();
    if (storedSessionIds) {
      console.info(`Found stored session: ${JSON.stringify(storedSessionIds)}`);
      // this.connect(storedSessionIds.sessionToken, storedSessionIds.userId);
    }
  }

  getState() {
    return get(this.store);
  }

  startGame() {
    this.socket?.emit("StartGame");
  }

  hasSessionData(): boolean {
    const state = this.getState();
    if (state.sessionToken && state.userId) return true;
    const sessionIds = localStorage.getItem('currentSessionIds');
    if (!sessionIds) return false;
    const sessionIdsJson = JSON.parse(sessionIds);
    return typeof sessionIdsJson.userId === "string" && typeof sessionIdsJson.sessionToken === "string";
  }

  private checkPermissionBit(permissionNumber: number, bitIndex: number): boolean {
    return (permissionNumber & (1 << bitIndex)) > 0;
  }

  getPlayerPermissions(PlayerId?: string): PlayerPermissionObj {
    if (!PlayerId) PlayerId = this.getState().userId ?? undefined;
    const playerPermissionNumber: number = this.getState().players?.find((player) => player.PlayerId == PlayerId)?.Permissions ?? 0;
    return {
      isHost: this.checkPermissionBit(playerPermissionNumber, 0)
    };
  }

  subscribe = this.store.subscribe;

  private getStoredSessionIds(): { sessionToken: string, userId: string } | null {
    if (typeof window === 'undefined') return null;
    const sessionIds = localStorage.getItem('currentSessionIds');
    if (!sessionIds) return null;
    const sessionIdsJson = JSON.parse(sessionIds);
    if (typeof sessionIdsJson.userId !== "string" || typeof sessionIdsJson.sessionToken !== "string") {
      return null;
    }
    return { sessionToken: sessionIdsJson.sessionToken, userId: sessionIdsJson.userId };
  }

  private saveSessionIds(sessionToken: string, userId: string) {
    if (typeof window !== 'undefined') {
      localStorage.setItem('currentSessionIds', JSON.stringify({ sessionToken, userId, joinCode: this.getState().joinCode }));
    }
  }

  private saveJoinCode() {
    const sessionIds = localStorage.getItem('currentSessionIds');
    if (!sessionIds) return;
    const sessionIdsJson = JSON.parse(sessionIds);
    localStorage.setItem('currentSessionIds', JSON.stringify({ sessionToken: sessionIdsJson.sessionToken, userId: sessionIdsJson.userId, joinCode: this.getState().joinCode }));
  }

  private clearSessionIds() {
    if (typeof window !== 'undefined') {
      const sessionIds = localStorage.getItem('currentSessionIds');
      if (!sessionIds) return;
      const sessionIdsJson = JSON.parse(sessionIds);
      const lastSessionData = { joinCode: sessionIdsJson.joinCode };
      localStorage.setItem('lastSessionIds', JSON.stringify(lastSessionData));
      localStorage.removeItem('currentSessionIds');
    }
  }

  isConnected(): boolean {
    return this.socket?.connected ?? false;
  }

  hasRoomData(): boolean {
    return get(this.store).gameState != -1;
  }

  getUserId(): string | undefined {
    return this.getState().userId ?? undefined;
  }

  getUser(playerId?: string): PlayerObj | undefined {
    if (!playerId) playerId = this.getUserId();
    return this.getState().players.find((player) => player.PlayerId == playerId);
  }

  kickPlayer(playerId: string) {
    if (!this.getPlayerPermissions().isHost) return;
    this.socket?.emit("KickPlayer", JSON.stringify({ PlayerId: playerId }));
  }

  renamePlayer(playerId: string | undefined, newName: string) {
    if (!playerId) playerId = this.getUserId();
    if (!this.getPlayerPermissions().isHost && playerId != this.getUserId()) return;
    this.socket?.emit("UpdatePlayer", JSON.stringify({ PlayerId: playerId, Username: newName }));
  }

  isCurrentPlayer(playerId: string): boolean {
    return this.getState().userId == playerId;
  }

  connect(sessionToken?: string, userId?: string) {
    if (!sessionToken) sessionToken = this.getState().sessionToken || undefined;
    if (!userId) userId = this.getState().userId || undefined;

    if (!sessionToken || !userId) {
      const storedSessionIds = this.getStoredSessionIds();
      if (!sessionToken) sessionToken = storedSessionIds?.sessionToken;
      if (!userId) userId = storedSessionIds?.userId;
    }

    if (this.socket) {
      console.warn(`Socket already connected! Rejecting new connection to ${sessionToken}`);
      return;
    }

    this.socket = io({
      transports: ['websocket'],
      query: { sessionToken },
    });

    this.setupSocketEventHandlers(sessionToken, userId);
  }

  private setupSocketEventHandlers(sessionToken: string, userId: string) {
    this.socket?.on('connect', () => this.handleConnect(sessionToken, userId));
    this.socket?.on('disconnect', this.handleDisconnect.bind(this));
    this.socket?.on('Status', this.handleStatus.bind(this));
    this.socket?.on('RoomInfo', this.handleRoomInfo.bind(this));
    this.socket?.on('error', this.handleError.bind(this));
  }

  private handleConnect(sessionToken: string, userId: string) {
    console.info('Connected to room');
    this.saveSessionIds(sessionToken, userId);
    window.history.replaceState({}, "", "/Game");
    this.store.update((state) => ({
      ...state,
      socket: this.socket,
      userId,
      connected: true,
      sessionToken,
    }));
  }

  private handleDisconnect() {
    console.info('Disconnected from server');
    this.store.update((state) => ({ ...state, connected: false }));
  }

  private handleStatus(message: StatusInfoObj) {
    console.log("Status: ", message);
    if (message.IsError && message.StatusCode !== "connection_from_different_socket") {
      this.leaveRoom();
    }
    if (message.IsError) {
      this.socket = null;
      window.history.replaceState({}, "", "/");
    }
    this.store.update((state) => ({
      ...state,
      messages: [...state.messages, message],
    }));
  }

  private handleRoomInfo(message: RoomInfoObj) {
    console.log("RoomInfo: ", message);
    this.store.update((state) => ({
      ...state,
      roomId: message.RoomId,
      joinCode: message.JoinCode,
      gameState: message.GameState,
      cardDeckId: message.CardDeckId,
      players: message.Players,
    }));
    this.saveJoinCode();
    this.store.update((state) => ({
      ...state,
      messages: [...state.messages, message],
    }));
  }

  private handleError(error: string) {
    console.error('Socket error:', error);
  }

  sendMessage(message: string) {
    if (this.socket && message.trim()) {
      this.socket.emit('event', message);
    }
  }

  leaveRoom() {
    console.log("leave room");
    if (this.socket) {
      this.socket.disconnect();
      this.socket = null;
    }
    if (this.getState().sessionToken) {
      fetch(`/api/room/leave`, {
        method: "POST",
        body: JSON.stringify({
          SessionToken: this.getState().sessionToken
        }),
        headers: {
          "Content-Type": "application/json",
        },
      });
    }
    this.clearSessionIds();
    this.store.set({
      roomId: null,
      joinCode: null,
      gameState: -1,
      gameOptions: {},
      players: [],
      cardDeckId: null,
      socket: null,
      connected: false,
      userId: null,
      messages: [],
      sessionToken: null,
    });
    window.history.replaceState({}, "", "/");
  }
}

export const sessionStore = new SessionManager();
