import ClassicCard from "../components/Cards/ClassicCard.svelte";
import HexV1Card from "../components/Cards/HexV1Card.svelte";

interface CardDeck {
    id: number;
    name: string;
    cardComponent: any;
}

export const CardDecks: CardDeck[] = [
    { id: 0, name: "Classic", cardComponent: ClassicCard },
    { id: 1, name: "HexV1", cardComponent: HexV1Card },
];
