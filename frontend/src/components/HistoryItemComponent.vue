<script lang="ts" setup>
import {history} from "../../wailsjs/go/models";
import {reactive} from "vue";
import {Rollback} from "../../wailsjs/go/app/App";

interface Props {
  game: history.GameHistory
}

const props = defineProps<Props>()

interface Data {
  toggled: boolean,
}

const data = reactive<Data>({
  toggled: false
})

function toggleStats() {
  data.toggled = !data.toggled
}

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
  <tr v-on:click="toggleStats">
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
<!--  <div class="modal fade" id="extended-history-{{props.game.gameId}}" tabindex="-1" aria-labelledby="label-extended-history-{{props.game.gameId}}" aria-hidden="true">-->
<!--    <div class="modal-dialog modal-dialog modal-dialog-centered bg-transparent text-black">-->
<!--      <div class="modal-content">-->
<!--        <div class="modal-header">-->
<!--          <h5 class="modal-title" id="label-extended-history-{{props.game.gameId}}">Modal title</h5>-->
<!--          <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>-->
<!--        </div>-->
<!--        <div class="modal-body">-->
<!--          <div class="container">-->
<!--            <div class="row">-->
<!--              <div class="col-6">-->
<!--                <label for="exampleFormControlInput1" class="form-label">Email address</label>-->
<!--                <input type="text" class="form-control" id="exampleFormControlInput1" placeholder="name@example.com">-->
<!--              </div>-->
<!--              <div class="col-6"></div>-->
<!--            </div>-->
<!--          </div>-->
<!--        </div>-->
<!--        <div class="modal-footer">-->
<!--          <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>-->
<!--        </div>-->
<!--      </div>-->
<!--    </div>-->
<!--  </div>-->
</template>
