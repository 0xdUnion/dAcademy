<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import axios from 'axios'
import { marked } from 'marked'

const route = useRoute()
const courseSlug = route.params.courseSlug
const chapterID = route.params.chapterID

const respond = ref({})
const sections = ref([])
const current = ref(0)
const loading = ref(true)
const finished = ref(false)

// quiz
const quizData = ref([])
const quizActive = ref(false)
const currentQuizList = ref([])
const quizIndex = ref(0)
const currentQuiz = ref(null)
const userAnswer = ref('')
const quizResult = ref(null)
const showAnswers = ref(false)

onMounted(async () => {
  try {
    const res = await axios.get(`/api/chapter/${courseSlug}/${chapterID}`)
    respond.value = res.data || {}
    sections.value = res.data.sections || []
    quizData.value = res.data.quiz || []
  } catch (err) {
    console.error('Failed to load chapter:', err)
  } finally {
    loading.value = false
  }
})

function buildQuizListForCurrentSection() {
  const qids = sections.value[current.value].quiz
  if (!Array.isArray(qids) || qids.length === 0) {
    currentQuizList.value = []
    return
  }
  currentQuizList.value = qids
    .map((id) => quizData.value.find((q) => q.id === id))
    .filter(Boolean)
}

function startQuiz() {
  buildQuizListForCurrentSection()
  if (!currentQuizList.value.length) return
  quizIndex.value = 0
  currentQuiz.value = currentQuizList.value[quizIndex.value]
  quizActive.value = true
  quizResult.value = null
  userAnswer.value = ''
  showAnswers.value = false
}

function loadNextQuizInList() {
  if (quizIndex.value < currentQuizList.value.length - 1) {
    quizIndex.value++
    currentQuiz.value = currentQuizList.value[quizIndex.value]
    userAnswer.value = ''
    quizResult.value = null
    showAnswers.value = false
  } else {
    quizActive.value = false
    setTimeout(() => nextSection(), 200)
  }
}

function checkQuiz() {
  if (!currentQuiz.value) return
  const answers = (currentQuiz.value.answer || []).map((a) =>
    String(a).trim().toLowerCase()
  )
  const given = String(userAnswer.value || '').trim().toLowerCase()
  if (answers.includes(given)) {
    quizResult.value = '✅ Correct!'
    setTimeout(() => {
      loadNextQuizInList()
    }, 600)
  } else {
    quizResult.value = '❌ Try again!'
  }
}

function revealAnswer() {
  showAnswers.value = !showAnswers.value
}

function nextSection() {
  if (current.value < sections.value.length - 1) {
    current.value++
  } else {
    finished.value = true
  }
  quizActive.value = false
  currentQuizList.value = []
  currentQuiz.value = null
  quizIndex.value = 0
  userAnswer.value = ''
  quizResult.value = null
  showAnswers.value = false
}
</script>

<template>
  <div>
    <h2>📘 Chapter {{ respond.chapter_id }} : {{ respond.chapter_title }}</h2>

    <div v-if="loading">Loading...</div>

    <div v-else>
      <div v-if="!finished">
        <div v-if="sections.length">
          <div v-if="sections[current].type === 'text'">
            <p>{{ sections[current].text }}</p>
          </div>

          <div v-else-if="sections[current].type === 'markdown'">
            <div v-html="marked.parse(sections[current].text)"></div>
          </div>

          <!-- Quiz UI -->
          <div v-if="quizActive && currentQuiz">
            <h3>🧠 Quiz <small>({{ quizIndex + 1 }} / {{ currentQuizList.length }})</small></h3>
            <p>{{ currentQuiz.text }}</p>

            <div v-if="currentQuiz.type === 'input'">
              <input
                v-model="userAnswer"
                placeholder="Your answer..."
                @keyup.enter="checkQuiz"
                class="border p-1 rounded"
              />
              <button @click="checkQuiz">Submit</button>
              <button @click="revealAnswer">
                {{ showAnswers ? 'Hide Answers' : '💡 Show Answer' }}
              </button>
            </div>

            <div v-else>
              <p>Unsupported quiz type: {{ currentQuiz.type }}</p>
              <button @click="checkQuiz">Mark as done</button>
            </div>

            <div v-if="showAnswers" class="answers">
              <p>📝 Accepted answers:</p>
              <ul>
                <li v-for="(ans, i) in currentQuiz.answer" :key="i">{{ ans }}</li>
              </ul>
            </div>

            <p v-if="quizResult">{{ quizResult }}</p>
          </div>

          <!-- Controls -->
          <div v-else>
            <button
              v-if="sections[current].quiz && sections[current].quiz.length"
              @click="startQuiz"
            >
              Start Quiz 🧩
            </button>

            <button
              v-else
              @click="nextSection"
            >
              Next ➡️
            </button>
          </div>
        </div>
        <div v-else>
          <p>No sections found.</p>
        </div>
      </div>

      <div v-else>
        <h3>🎉 Chapter Finished!</h3>
        <p>Well done, learner!</p>
        <RouterLink :to="`/c/${courseSlug}`">Back</RouterLink>
      </div>
    </div>
  </div>
</template>

<style scoped>
button {
  margin-top: 8px;
  margin-right: 6px;
  padding: 6px 12px;
  border-radius: 6px;
  border: none;
  background-color: #222;
  color: #fff;
  cursor: pointer;
}
button:hover {
  background-color: #444;
}
.answers {
  margin-top: 10px;
  background: #111;
  padding: 10px;
  border-radius: 6px;
  color: white;
}
.answers ul {
  margin: 0;
  padding-left: 20px;
}
</style>
