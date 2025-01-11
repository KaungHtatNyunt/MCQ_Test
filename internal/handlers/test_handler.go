package handlers

import (
	"html/template"
	"math/rand"
	"mcq-test-system/internal/models"
	"mcq-test-system/internal/utils"
	"net/http"
	"strconv"
	"time"
)

type TestHandler struct {
	templates *template.Template
	sessions  map[string]*models.TestSession
}

func NewTestHandler() *TestHandler {
	return &TestHandler{
		templates: template.Must(template.ParseGlob("/opt/render/project/go/src/github.com/KaungHtatNyunt/MCQ_Test/templates/*.html")),
		sessions:  make(map[string]*models.TestSession),
	}
}

// Test စမယ်။ StartTest နဲ့ စမယ်။ ဒါတွေ စလုပ်မယ်။
//
//	GET request လုပ်ရင် start.html ကို ExecuteTemplate လုပ်မယ်။ "start.html" ပေါ်အောင် လုပ်ပြီးသား
//
// Excel ဖိုင် ဖတ်မယ်။ အဖြေကို ရောမယ်။ new session  ကို Create လုပ်မယ်
// ကွက်ကီး ထည့်မယ်။ http.Redirect(w, r, "/question",
func (h *TestHandler) StartTest(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		h.templates.ExecuteTemplate(w, "start.html", nil) //
		return
	}

	numQuestions, _ := strconv.Atoi(r.FormValue("num_questions"))
	if numQuestions < 1 {
		http.Error(w, "Invalid number of questions", http.StatusBadRequest)
		return
	}

	questions, err := utils.LoadQuestionsFromExcel("/opt/render/project/go/src/github.com/KaungHtatNyunt/MCQ_Test/MCQ_question.xlsx")
	if err != nil {
		http.Error(w, "Failed to load questions", http.StatusInternalServerError)
		return
	}
	
	// Shuffle questions
	rand.Shuffle(len(questions), func(i, j int) {
		questions[i], questions[j] = questions[j], questions[i]
	})

	// Create new session
	sessionID := strconv.FormatInt(time.Now().UnixNano(), 10)
	session := &models.TestSession{
		Questions:    questions[:numQuestions],
		StartTime:    time.Now(),
		NumQuestions: numQuestions,
		Answers:      make([]int, numQuestions),
	}
	h.sessions[sessionID] = session

	http.SetCookie(w, &http.Cookie{
		Name:  "session_id",
		Value: sessionID,
		Path:  "/",
	})

	http.Redirect(w, r, "/question", http.StatusSeeOther)
}

func (h *TestHandler) HandleQuestion(w http.ResponseWriter, r *http.Request) {
	session := h.getSession(r)
	if session == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Check if test time exceeded
	elapsed := time.Since(session.StartTime)
	allowedTime := time.Duration(session.NumQuestions) * 2 * time.Minute
	if elapsed > allowedTime {
		session.EndTime = time.Now()
		http.Redirect(w, r, "/report", http.StatusSeeOther)
		return
	}

	if session.CurrentIndex >= session.NumQuestions {
		session.EndTime = time.Now()
		http.Redirect(w, r, "/report", http.StatusSeeOther)
		return
	}

	currentQ := session.Questions[session.CurrentIndex]
	// Shuffle answers
	answers := make([]string, len(currentQ.Answers))
	copy(answers, currentQ.Answers)
	rand.Shuffle(len(answers), func(i, j int) {
		answers[i], answers[j] = answers[j], answers[i]
	})

	data := struct {
		Question    models.Question
		Answers     []string
		QuestionNum int
		TimeLeft    int64
	}{
		Question:    currentQ,
		Answers:     answers,
		QuestionNum: session.CurrentIndex + 1,
		TimeLeft:    int64(allowedTime-elapsed) / int64(time.Second),
	}

	h.templates.ExecuteTemplate(w, "question.html", data)
}

func (h *TestHandler) HandleSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session := h.getSession(r)
	if session == nil {
		http.Error(w, "Session not found", http.StatusBadRequest)
		return
	}

	answer, _ := strconv.Atoi(r.FormValue("answer"))
	correct := answer == session.Questions[session.CurrentIndex].CorrectAnswer
	if correct {
		session.CorrectCount++
	}
	session.Answers[session.CurrentIndex] = answer
	session.CurrentIndex++

	w.Write([]byte(`{"correct": ` + strconv.FormatBool(correct) + `}`))
}

func (h *TestHandler) GenerateReport(w http.ResponseWriter, r *http.Request) {
	session := h.getSession(r)
	if session == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if session.EndTime.IsZero() {
		session.EndTime = time.Now()
	}

	timeTaken := session.EndTime.Sub(session.StartTime)
	avgTimePerQ := float64(timeTaken.Minutes()) / float64(session.NumQuestions)
	scorePercentage := float64(session.CorrectCount) / float64(session.NumQuestions) * 100

	report := models.Report{
		NumQuestions:    session.NumQuestions,
		Topic:           session.Questions[0].Topic,
		CorrectAnswers:  session.CorrectCount,
		ScorePercentage: scorePercentage,
		TimeTaken:       timeTaken,
		AvgTimePerQ:     avgTimePerQ,
		StartTime:       session.StartTime,
		EndTime:         session.EndTime,
	}

	h.templates.ExecuteTemplate(w, "report.html", report)
}

func (h *TestHandler) getSession(r *http.Request) *models.TestSession {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return nil
	}
	return h.sessions[cookie.Value]
}
