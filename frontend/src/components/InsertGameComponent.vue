<script lang="ts" setup>
import {reactive} from 'vue'
import {Insert} from '../../wailsjs/go/app/App'
import {history, player} from "../../wailsjs/go/models";

interface Props {
  selectedPlayer: player.Player
}

const props = defineProps<Props>()

interface Data {
  insertionInProgress: boolean
  isFormValid: boolean
  input: history.InsertRequest
  errors: ValidationErrors
  suggestedElo?: number
}

interface ValidationErrors {
  totalGames?: string
  wins?: string
  gamesLeftEarly?: string
  highestWinStreak?: string
  winsStreak?: string
  mvp?: string
  elo?: string
}

const averageEloPerWin = 15
const averageEloLossPerLose = 3
const averageEloPerLeave = 13
const defaultElo = 1500

const data = reactive<Data>({
  insertionInProgress: false,
  isFormValid: true,
  input: new history.InsertRequest({
    totalGames: 1,
    wins: 1,
    gamesLeftEarly: 0,
    highestWinStreak: 1,
    winsStreak: 1,
    mvp: 1,
    elo: defaultElo + averageEloPerWin,
    timestamp: (new Date()).setDate((new Date()).getDate()-1)
  }),
  errors: {
    totalGames: undefined,
    wins: undefined,
    gamesLeftEarly: undefined,
    highestWinStreak: undefined,
    winsStreak: undefined,
    mvp: undefined,
    elo: undefined,
  },
  suggestedElo: defaultElo + averageEloPerWin
})

const validateForm = () => {
  validation.totalGames()
  validation.wins()
  validation.gamesLeftEarly()
  validation.highestWinStreak()
  validation.winsStreak()
  validation.mvp()
  validation.elo()

  data.suggestedElo = calculateSuggestedElo()

  data.isFormValid = data.errors.totalGames == undefined
      && data.errors.wins == undefined
      && data.errors.gamesLeftEarly == undefined
      && data.errors.highestWinStreak == undefined
      && data.errors.winsStreak == undefined
      && data.errors.mvp == undefined
      && data.errors.elo == undefined
}

const validation = {
  totalGames: () => {
    if (data.input.totalGames < 1) {
      data.errors.totalGames = 'Total Games cannot be less than 1.';
      return
    }
    if (data.input.totalGames < (data.input.wins + data.input.gamesLeftEarly)) {
      data.errors.totalGames = `Total Games must be at least ${data.input.wins + data.input.gamesLeftEarly} sum of Wins and Games Left Early.`;
      return
    }
    data.errors.totalGames = undefined
  },
  wins: () => {
    if (data.input.wins < 0) {
      data.errors.wins = 'Wins cannot be less than 0.';
      return
    }
    if (data.input.wins > (data.input.totalGames - data.input.gamesLeftEarly)) {
      data.errors.wins = `Wins cannot be higher than ${data.input.totalGames - data.input.gamesLeftEarly} Total Games - Games Left Early.`;
      return;
    }
    data.errors.wins = undefined
  },
  gamesLeftEarly: () => {
    if (data.input.gamesLeftEarly < 0) {
      data.errors.gamesLeftEarly = 'Games Left Early cannot be less than 0.';
      return
    }
    if (data.input.gamesLeftEarly > (data.input.totalGames - data.input.wins)) {
      data.errors.gamesLeftEarly = `Games Left Early cannot be higher than ${data.input.totalGames - data.input.wins} Total Games - Wins.`;
      return;
    }
    data.errors.gamesLeftEarly = undefined;
  },
  highestWinStreak: () => {
    if (data.input.highestWinStreak > data.input.wins) {
      data.errors.highestWinStreak = `Highest Wins Streak cannot be higher than ${data.input.wins} Wins.`;
      return;
    }
    if (data.input.highestWinStreak != data.input.winsStreak && data.input.highestWinStreak > data.input.wins - data.input.winsStreak) {
      data.errors.highestWinStreak = `Highest Wins Streak should be ${data.input.winsStreak}, or not higher than ${data.input.wins - data.input.winsStreak} Wins - Current Wins Streak.`;
      return;
    }
    data.errors.highestWinStreak = undefined;
  },
  winsStreak: () => {
    if (data.input.winsStreak > data.input.wins) {
      data.errors.winsStreak = `Current Wins Streak cannot be higher than ${data.input.wins} Wins.`;
      return;
    }
    if (data.input.winsStreak != data.input.highestWinStreak && data.input.winsStreak > data.input.wins - data.input.highestWinStreak) {
      data.errors.winsStreak = `Current Wins Streak should be ${data.input.highestWinStreak}, or not higher than ${data.input.wins - data.input.highestWinStreak} Wins - Highest Wins Streak.`;
      return;
    }
    data.errors.winsStreak = undefined;
  },
  mvp: () => {
    if (data.input.mvp > data.input.totalGames - data.input.gamesLeftEarly) {
      data.errors.mvp = `MVP cannot be higher than ${data.input.totalGames - data.input.gamesLeftEarly} Total Games - Games Left Early.`;
      return;
    }
    data.errors.mvp = undefined;
  },
  elo: () => {
    if (data.input.elo < 1000) {
      data.errors.elo = `ELO cannot be less than 1000.`;
      return;
    }
    if (data.input.elo > 3000) {
      data.errors.elo = `ELO cannot be more than 3000.`;
      return;
    }
    data.errors.elo = undefined;
  }
}

const insertGame = (event: Event) => {
  event.preventDefault()
  validateForm()
  if (!data.isFormValid) {
    return
  }

  data.insertionInProgress = true
  Insert(props.selectedPlayer, data.input).then(() => {
    data.input = new history.InsertRequest()
  }).finally(() => {
    data.insertionInProgress = false
  })
}
const calculateSuggestedElo = (): number => {
  const eloGained = data.input.wins * averageEloPerWin
  const eloLostFromLeftGames = data.input.gamesLeftEarly * averageEloPerLeave
  const eloLostFromLosses = (data.input.totalGames - data.input.wins - data.input.gamesLeftEarly) * averageEloLossPerLose

  return (defaultElo + eloGained) - (eloLostFromLosses + eloLostFromLeftGames)
}

</script>

<template>
  <form class="row g-3" v-on:submit="insertGame">
    <div class="col-md-6">
      <label for="totalGames" class="form-label">Total Games</label>
      <input id="totalGames" type="number"
             v-model="data.input.totalGames"
             @input="validateForm"
             :class="{'is-invalid': data.errors.totalGames}"
             class="form-control">
      <div v-if="data.errors.totalGames" class="invalid-feedback">
        {{ data.errors.totalGames }}
      </div>
    </div>
    <div class="col-md-6">
      <label for="wins" class="form-label">Wins</label>
      <input id="wins" type="number" min="0"
             v-model="data.input.wins"
             @input="validateForm"
             :class="{'is-invalid': data.errors.wins}"
             class="form-control">
      <div v-if="data.errors.wins" class="invalid-feedback">
        {{ data.errors.wins }}
      </div>
    </div>
    <div class="col-md-6">
      <label for="gameLeftEarly" class="form-label">Games Left Early</label>
      <input id="gameLeftEarly" type="number" min="0"
             v-model="data.input.gamesLeftEarly"
             :class="{'is-invalid': data.errors.gamesLeftEarly}"
             @input="validateForm"
             class="form-control">
      <div v-if="data.errors.gamesLeftEarly" class="invalid-feedback">
        {{ data.errors.gamesLeftEarly }}
      </div>
    </div>
    <div class="col-md-6">
      <label for="gamesLost" class="form-label">Games Lost (Calculated)</label>
      <input id="gamesLost" type="number"
             v-bind:value="data.input.totalGames-data.input.wins-data.input.gamesLeftEarly"
             class="form-control" disabled>
      <p class="fst-italic text-secondary small fw-bold mt-1">Calculated based on Total Games - (Wins + Games Left
        Early).</p>
    </div>
    <div class="col-md-6">
      <label for="highestWinsStreak" class="form-label">Highest Wins Streak</label>
      <input id="highestWinsStreak" type="number"
             v-model="data.input.highestWinStreak"
             @input="validateForm"
             :class="{'is-invalid': data.errors.highestWinStreak}"
             class="form-control">
      <div v-if="data.errors.highestWinStreak" class="invalid-feedback">
        {{ data.errors.highestWinStreak }}
      </div>
    </div>
    <div class="col-md-6">
      <label for="winsStreak" class="form-label">Current Wins Streak</label>
      <input id="winsStreak" type="number"
             v-model="data.input.winsStreak"
             @input="validateForm"
             :class="{'is-invalid': data.errors.winsStreak}"
             class="form-control">
      <div v-if="data.errors.winsStreak" class="invalid-feedback">
        {{ data.errors.winsStreak }}
      </div>
    </div>
    <div class="col-md-6">
      <label for="mvp" class="form-label">MVP</label>
      <input id="mvp" type="number" min="0"
             v-model="data.input.mvp"
             @input="validateForm"
             :class="{'is-invalid': data.errors.mvp}"
             data-bs-toggle="tooltip" data-bs-placement="top"
             class="form-control">
      <div v-if="data.errors.mvp" class="invalid-feedback">
        {{ data.errors.mvp }}
      </div>
    </div>
    <div class="col-md-6">
      <label for="timestamp" class="form-label">Date <span class="text-secondary fw-bold">(Month/Day/Year, Hour:Minute)</span></label>
      <VueDatePicker
          v-model="data.input.timestamp" :clearable="false"
          :min-date="(new Date()).setFullYear((new Date()).getFullYear()-1)"
          :max-date="new Date()"
          :enable-minutes="true"
      />
      <div v-if="data.errors.gamesLeftEarly" class="invalid-feedback">
        {{ data.errors.gamesLeftEarly }}
      </div>
    </div>
    <div class="col-md-6">
      <label for="elo" class="form-label">ELO</label>
      <input id="elo" type="number" min="0"
             v-model="data.input.elo"
             @input="validateForm"
             :class="{'is-invalid': data.errors.elo}"
             data-bs-toggle="tooltip" data-bs-placement="top"
             class="form-control">
      <div v-if="data.errors.elo" class="invalid-feedback">
        {{ data.errors.elo }}
      </div>
      <p class="fst-italic text-secondary small fw-bold mt-1">The suggested ELO calculated based on your inputs is:
        <span class="text-success">{{ data.suggestedElo }}</span>.<br>You could set it higher, but it might be <span
            class="text-decoration-underline">suspicious</span>.</p>
    </div>
    <div class="col-12">
      <button type="submit" class="btn btn-primary" v-bind:disabled="data.insertionInProgress || !data.isFormValid">
        Insert
      </button>
    </div>
  </form>
</template>
