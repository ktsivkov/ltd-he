<script lang="ts" setup>
import {history} from "../../wailsjs/go/models";
import {Rollback} from "../../wailsjs/go/app/App";

interface Props {
  game: history.GameHistory
}

const props = defineProps<Props>()

function getClassByOutcome() {
  switch (props.game.outcome){
    case "WIN":
      return "bg-success"
    case "LOSS":
      return "bg-danger"
    case "LEAVE":
      return "bg-warning"
    case "DRAW":
      return "bg-secondary"
  }
}

function rollback() {
  Rollback(props.game).catch(error => console.error(error))
}

</script>

<template>
  <tr v-if="props.game">
    <td class="text-center">
      <span class="badge badge rounded-pill" :class="getClassByOutcome()">{{props.game.outcome}}</span>
    </td>
    <td>
      {{ props.game.elo }}
      <span v-if="['WIN', 'LOSS', 'LEAVE'].includes(props.game.outcome)" class="badge badge rounded-pill"
            :class="getClassByOutcome()">
      {{ `${props.game.outcome == 'WIN' ? '+' : '-'}${props.game.eloDiff}` }}
    </span>
      <span v-if="props.game.winsStreak > 1">&#128293; {{ props.game.winsStreak }}</span>
    </td>
    <td>{{ props.game.date }}</td>
    <td>
      <button type="button" class="btn btn-outline-success" v-if="!props.game.isLast" v-on:click="rollback">Restore</button>
    </td>
  </tr>
</template>
