// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {player} from '../models';
import {history} from '../models';

export function BackupFolder(arg1:player.Player):Promise<string>;

export function Insert(arg1:player.Player,arg2:history.InsertRequest):Promise<void>;

export function ListPlayers():Promise<Array<player.Player>>;

export function LoadHistory(arg1:player.Player):Promise<history.History>;

export function Rollback(arg1:history.GameHistory):Promise<void>;
