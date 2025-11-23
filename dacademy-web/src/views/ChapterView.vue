<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import Nav from '@/views/components/Nav.vue'
import axios from 'axios'
import { Marked } from 'marked'
import { markedHighlight } from "marked-highlight";
import hljs from 'highlight.js';
import 'highlight.js/styles/github-dark-dimmed.css'

const route = useRoute()
const marked = new Marked(
  markedHighlight({
    emptyLangClass: 'hljs',
    langPrefix: 'hljs language-',
    highlight(code, lang, info) {
      const language = hljs.getLanguage(lang) ? lang : 'plaintext';
      return hljs.highlight(code, { language }).value;
    }
  })
);

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

function skipQuiz() {
  quizResult.value = '➡️ Skipped'
  setTimeout(() => {
    loadNextQuizInList()
  }, 600)
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
  <Nav />
  <div class="container mx-auto p-6 space-y-6">
    <h2 class="text-2xl font-semibold text-primary">
      📘 Chapter {{ respond.chapter_id }} : {{ respond.chapter_title }}
    </h2>

    <div v-if="loading" class="text-center text-lg text-gray-500">Loading...</div>

    <div v-else>
      <div v-if="!finished">
        <div v-if="sections.length">
          <div v-if="sections[current].type === 'text'">
            <p class="text-lg text-gray-700">{{ sections[current].text }}</p>
          </div>

          <div v-else-if="sections[current].type === 'markdown'">
            <div v-html="marked.parse(sections[current].text)"></div>
          </div>

          <!-- Quiz UI -->
          <div v-if="quizActive && currentQuiz" class="mt-6">
            <h3 class="text-xl font-semibold text-primary">
              🧠 Quiz <small>({{ quizIndex + 1 }} / {{ currentQuizList.length }})</small>
            </h3>
            <p class="text-lg mt-2">{{ currentQuiz.text }}</p>

            <!-- Input -->
            <div v-if="currentQuiz.type === 'input'" class="mt-4 space-y-4">
              <input
                v-model="userAnswer"
                placeholder="Your answer..."
                @keyup.enter="checkQuiz"
                class="input input-bordered w-full"
              />
              <div class="space-x-4">
                <button @click="checkQuiz" class="btn btn-accent">Submit</button>
                <button @click="revealAnswer" class="btn btn-ghost">
                  {{ showAnswers ? 'Hide Answers' : '💡 Show Answer' }}
                </button>
              </div>
            </div>

            <!-- Single Choose -->
            <div v-else-if="currentQuiz.type === 'singleChoose'" class="mt-4 space-y-4">
              <div class="space-x-4">
                <button
                  v-for="(opt, i) in currentQuiz.options"
                  :key="i"
                  @click="userAnswer = opt"
                  :class="{'btn btn-primary': userAnswer === opt, 'btn': userAnswer !== opt}"
                >
                  {{ opt }}
                </button>
              </div>
              <div class="space-x-4">
                <button @click="checkQuiz" class="btn btn-accent">Submit</button>
                <button @click="revealAnswer" class="btn btn-ghost">
                  {{ showAnswers ? 'Hide Answers' : '💡 Show Answer' }}
                </button>
              </div>
            </div>

            <div v-else>
              <p>Unsupported quiz type: {{ currentQuiz.type }}</p>
              <button @click="skipQuiz" class="btn btn-warning">Mark as done</button>
            </div>

            <div v-if="showAnswers" class="answers bg-base-200 p-4 rounded-md mt-4">
              <p>📝 Accepted answers:</p>
              <ul>
                <li v-for="(ans, i) in currentQuiz.answer" :key="i">{{ ans }}</li>
              </ul>
            </div>

            <p v-if="quizResult" class="text-xl font-semibold mt-4">{{ quizResult }}</p>
          </div>

          <!-- Controls -->
          <div v-else class="mt-6">
            <button
              v-if="sections[current].quiz && sections[current].quiz.length"
              @click="startQuiz"
              class="btn btn-secondary"
            >
              Start Quiz 🧩
            </button>

            <button
              v-else
              @click="nextSection"
              class="btn btn-primary w-full"
            >
              Next ➡️
            </button>
          </div>
        </div>
        <div v-else>
          <p>No sections found.</p>
        </div>
      </div>

      <div v-else class="text-center">
        <h3 class="text-3xl font-bold text-green-500">🎉 Chapter Finished!</h3>
        <p class="text-lg">Well done, learner!</p>
        <RouterLink :to="`/c/${courseSlug}`" class="btn btn-accent mt-4">
          Back to Course
        </RouterLink>
      </div>
    </div>
  </div>
</template>

<style scoped>
.container {
  max-width: 1024px;
}

.answers {
  background: #1f2937;
  color: #fff;
  border-radius: 8px;
}
</style>
