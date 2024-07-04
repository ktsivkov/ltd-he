<script lang="ts" setup>
import {reactive} from 'vue'
import {Append} from '../../wailsjs/go/app/App'
import {history, player} from "../../wailsjs/go/models";

interface Props {
  selectedPlayer: player.Player
}

const props = defineProps<Props>()

interface Data {
  appendingInProgress: boolean
  input: history.AppendRequest
}

const data = reactive<Data>({
  appendingInProgress: false,
  input: new history.AppendRequest()
})

function appendGame() {
  data.appendingInProgress = true
  Append(props.selectedPlayer, data.input).then(()=>{data.input = new history.AppendRequest()}).finally(() => {
    data.appendingInProgress = false
  })
}

</script>

<template>
  <div class="row">
    <div class="col-6">
      <fieldset v-bind:disabled="data.appendingInProgress">
        <form class="row gy-2 gx-3 align-items-center">
          <div class="col">
            <label class="visually-hidden" for="elo">ELO</label>
            <input type="number" class="form-control" id="elo" placeholder="ELO" v-model="data.input.elo">
          </div>
          <div class="col-auto">
            <div class="form-check">
              <input class="form-check-input" type="checkbox" id="mvp" v-model="data.input.mvp">
              <label class="form-check-label" for="mvp">
                MVP
              </label>
            </div>
          </div>
          <div class="col-auto">
            <button type="submit" class="btn btn-outline-success" v-on:click="appendGame">Append Game</button>
          </div>
        </form>
      </fieldset>
    </div>
  </div>
</template>
