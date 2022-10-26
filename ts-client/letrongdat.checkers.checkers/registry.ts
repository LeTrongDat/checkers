import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgRejectGame } from "./types/checkers/tx";
import { MsgPlayMove } from "./types/checkers/tx";
import { MsgCreateGame } from "./types/checkers/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/letrongdat.checkers.checkers.MsgRejectGame", MsgRejectGame],
    ["/letrongdat.checkers.checkers.MsgPlayMove", MsgPlayMove],
    ["/letrongdat.checkers.checkers.MsgCreateGame", MsgCreateGame],
    
];

export { msgTypes }