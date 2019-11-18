package controllers

import (
	"Go-School/models"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type ExamController struct {
	CommonController
}

func (this *ExamController) ExamSubmit() {
	paramMap := this.checkParams([]string{"id", "answers", "time_lenght", "start_time"})

	// 时间戳
	datetime := time.Now().Unix()

	// 考试任务信息
	examRecordModel := this.findExamRecordInfoById(paramMap["id"].(string), "id", "exam_id", "already_repeat_count")

	// 考试信息
	examModel := this.findExamInfoById(examRecordModel.ExamId, "type", "integral", "exam_paper_id", "pass_score", "repeat_count", "endtime")

	// 验证能否再次考试
	this.checkAgainExam(examModel.RepeatCount, examRecordModel.AlreadyRepeatCount, examModel.Endtime, datetime)

	// 查询试卷信息
	examPaperModel := this.findExamPaperInfoById(examModel.ExamPaperId, "corpid", "is_reading")

	ormer := orm.NewOrm()
	// 判断当前是否补考
	var isRepeat int8
	if repeatCount, repeatCountErr := ormer.QueryTable(new(models.ExamRecordDetail)).Filter("corpid", AttributesStruct.Corpid).Filter("exam_record_id", paramMap["id"]).Count(); repeatCountErr != nil {
		this.ReturnData(-1, map[string]interface{}{})
	} else {
		if repeatCount > 0 {
			isRepeat = 1
		} else {
			isRepeat = 0
		}
	}

	var examRecordDetailModel models.ExamRecordDetail
	examRecordDetailModel.Corpid = AttributesStruct.Corpid
	examRecordDetailModel.ExamRecordId, _ = strconv.Atoi(paramMap["id"].(string))
	examRecordDetailModel.Start, _ = strconv.Atoi(paramMap["start_time"].(string))
	examRecordDetailModel.End = int(datetime)
	examRecordDetailModel.TimeCost, _ = strconv.Atoi(paramMap["time_lenght"].(string))
	examRecordDetailModel.IsRepeat = isRepeat
	examRecordDetailModel.Status = 0

	isExist := make(chan bool)

	go func() {
		// 开启事务
		ormer.Begin()

		// 创建考试记录
		examRecordDetailId, examRecordDetailErr := ormer.Insert(&examRecordDetailModel)
		if examRecordDetailErr != nil {
			// 事务回滚
			ormer.Rollback()
			this.ReturnData(59038, map[string]interface{}{})
		}

		// 记录最后一次考试记录时间
		examRecordModel.LastCreatetime = int(datetime)
		if examRecordNum, examRecordErr := ormer.Update(&examRecordModel, "last_createtime"); examRecordErr != nil || examRecordNum == 0 {
			// 事务回滚
			ormer.Rollback()
			this.ReturnData(59029, map[string]interface{}{})
		}

		// 检测用户答案
		answer, answerErr := this.checkAnswer(paramMap["answers"].([]interface{}), examModel.ExamPaperId, examPaperModel.Corpid, examRecordDetailId, examPaperModel.IsReading)
		if answerErr != nil {
			ormer.Rollback()
			this.ReturnData(59039, map[string]interface{}{})
		}

		// 记录用户答题记录
		if successNum, insertMultiErr := ormer.InsertMulti(1, answer); insertMultiErr != nil || successNum == 0 {
			ormer.Rollback()
			this.ReturnData(59040, map[string]interface{}{})
		}

		// 不需要阅卷
		if examPaperModel.IsReading == 0 {
			// 获取用户答题总分值
			var userScoreInt int
			for _, answ := range answer {
				score, _ := strconv.Atoi(answ.Score)
				userScoreInt += score
			}

			examRecordDetailModel.Id = int(examRecordDetailId)

			// 判断总分值是否达到及格线
			if userScoreInt >= examModel.PassScore { // 及格
				examRecordDetailModel.Score = strconv.Itoa(userScoreInt)
				examRecordDetailModel.Status = 1
				// 编辑考试记录详情 及格
				if upExamRecordDetailNum, upExamRecordDetailErr := ormer.Update(&examRecordDetailModel, "score", "status"); upExamRecordDetailErr != nil || upExamRecordDetailNum == 0 {
					ormer.Rollback()
					this.ReturnData(59028, map[string]interface{}{})
				}

				// 编辑分配的考试记录 已通过状态
				examRecordModel.LastScore = strconv.Itoa(userScoreInt)
				examRecordModel.Status = 4
				if upExamRecordNum, upExamRecordErr := ormer.Update(&examRecordModel, "last_score", "status"); upExamRecordErr != nil || upExamRecordNum == 0 {
					ormer.Rollback()
					this.ReturnData(59029, map[string]interface{}{})
				}

				// 考试通过加积分
				userModel := models.Users{Id: AttributesStruct.Userid}
				if userReadErr := ormer.Read(&userModel, "id", "score"); userReadErr != nil {
					ormer.Rollback()
					this.ReturnData(57004, map[string]interface{}{})
				} else {
					// 更改用户积分
					userModel.Score += examModel.Integral
					if upUser, upUserErr := ormer.Update(&userModel, "score"); upUserErr != nil || upUser == 0 {
						ormer.Rollback()
						this.ReturnData(57013, map[string]interface{}{})
					} else {
						// 积分更改成功，添加积分记录
						if _, addUserScoreErr := this.addUserScoreLog(examModel.Type, strconv.Itoa(examModel.Integral)); addUserScoreErr != nil {
							ormer.Rollback()
							this.ReturnData(57013, map[string]interface{}{})
						}
					}
				}
			} else { // 不及格
				examRecordDetailModel.Score = strconv.Itoa(userScoreInt)
				examRecordDetailModel.Status = 0
				// 编辑考试记录详情 及格
				if upExamRecordDetailNum, upExamRecordDetailErr := ormer.Update(&examRecordDetailModel, "score", "status"); upExamRecordDetailErr != nil || upExamRecordDetailNum == 0 {
					ormer.Rollback()
					this.ReturnData(59028, map[string]interface{}{})
				}

				// 判断考试是否支持补考
				if examModel.RepeatCount > 0 {
					// 判断是否还有补考次数
					if (examRecordModel.AlreadyRepeatCount + 1) < examModel.RepeatCount { // 还有补考次数
						// 编辑分配的考试记录 待补考状态
						examRecordModel.LastScore = strconv.Itoa(userScoreInt)
						examRecordModel.Status = 1
						if upExamRecordNum, upExamRecordErr := ormer.Update(&examRecordModel, "last_score", "status"); upExamRecordErr != nil || upExamRecordNum == 0 {
							ormer.Rollback()
							this.ReturnData(59029, map[string]interface{}{})
						}
					} else { // 没有补考次数
						examRecordModel.LastScore = strconv.Itoa(userScoreInt)
						examRecordModel.Status = 3
						if upExamRecordNum, upExamRecordErr := ormer.Update(&examRecordModel, "last_score", "status"); upExamRecordErr != nil || upExamRecordNum == 0 {
							ormer.Rollback()
							this.ReturnData(59029, map[string]interface{}{})
						}
					}
				} else { // 考试不支持补考

					// 考试未通过
					examRecordModel.LastScore = strconv.Itoa(userScoreInt)
					examRecordModel.Status = 3
					if upExamRecordNum, upExamRecordErr := ormer.Update(&examRecordModel, "last_score", "status"); upExamRecordErr != nil || upExamRecordNum == 0 {
						ormer.Rollback()
						this.ReturnData(59029, map[string]interface{}{})
					}
				}
			}
		} else { // 不需要阅卷
			examRecordModel.Status = 2 // 待阅卷
			if _, upExamRecordErr := ormer.Update(&examRecordModel, "status"); upExamRecordErr != nil {
				ormer.Rollback()
				this.ReturnData(59029, map[string]interface{}{})
			}
		}

		// 是补考的话，需要记录已补考次数
		if isRepeat == 1 {
			examRecordModel.AlreadyRepeatCount += 1
			if _, upExamRecordErr := ormer.Update(&examRecordModel, "already_repeat_count"); upExamRecordErr != nil {
				ormer.Rollback()
				this.ReturnData(59029, map[string]interface{}{})
			}
		}

		// 事务提交
		ormer.Commit()
		//ormer.Rollback()

		isExist <- true
	}()

	<-isExist

	var msg string
	if examPaperModel.IsReading == 1 {
		msg = "本次考试需要人工阅卷，请耐心等待"
	}

	returnMap := map[string]interface{}{
		"data": map[string]interface{}{
			"marking": examPaperModel.IsReading,
			"msg":     msg,
		},
	}

	this.ReturnData(0, returnMap)
}

// 考试任务信息
func (this *ExamController) findExamRecordInfoById(id string, fields ...string) models.ExamRecord {
	ormer := orm.NewOrm()

	var examRecord models.ExamRecord
	oneErr := ormer.QueryTable(new(models.ExamRecord)).Filter("id", id).Filter("corpid", AttributesStruct.Corpid).
		Filter("userid", AttributesStruct.Userid).One(&examRecord, fields...)
	if oneErr != nil {
		this.ReturnData(-1, map[string]interface{}{})
	}
	if oneErr == orm.ErrNoRows {
		this.ReturnData(59025, map[string]interface{}{})
	}

	return examRecord
}

// 考试信息
func (this *ExamController) findExamInfoById(id int, fields ...string) models.Exam {
	ormer := orm.NewOrm()

	var examModel models.Exam
	oneErr := ormer.QueryTable(new(models.Exam)).Filter("corpid", AttributesStruct.Corpid).Filter("id", id).One(&examModel, fields...)
	if oneErr != nil {
		this.ReturnData(-1, map[string]interface{}{})
	}

	if oneErr == orm.ErrNoRows {
		this.ReturnData(59023, map[string]interface{}{})
	}

	return examModel
}

// 验证能否再次考试
func (this *ExamController) checkAgainExam(repeatCount, alreadyRepeatCount int8, endtime int, datetime int64) bool {
	if datetime > int64(endtime) {
		this.ReturnData(59033, map[string]interface{}{})
	}

	if repeatCount > 0 && alreadyRepeatCount >= repeatCount {
		this.ReturnData(59037, map[string]interface{}{})
	}

	return true
}

// 查询试卷信息
func (this *ExamController) findExamPaperInfoById(id int, fields ...string) models.ExamPaper {
	ormer := orm.NewOrm()
	//cond := orm.NewCondition()

	var examPaperModel models.ExamPaper
	//oneErr := ormer.QueryTable(new(models.ExamPaper)).Filter("id", id).Filter("status", 1).SetCond(cond.AndCond(cond.And("corpid", AttributesStruct.Corpid).Or("corpid", beego.AppConfig.String("other::defaultCorpid")))).One(&examPaperModel, fields...)
	oneErr := ormer.QueryTable(new(models.ExamPaper)).Filter("id", id).Filter("status", 1).One(&examPaperModel, fields...)
	if oneErr != nil {
		this.ReturnData(-1, map[string]interface{}{})
	}

	if oneErr == orm.ErrNoRows {
		this.ReturnData(59018, map[string]interface{}{})
	}

	return examPaperModel
}

// 检测用户答案
func (this *ExamController) checkAnswer(answers []interface{}, examPaperId, corpid int, examRecordDetailId int64, isReading int8) ([]models.ExamRecordTestDetail, error) {
	ormer := orm.NewOrm()
	var examPaperTestsModel []models.ExamPaperTest
	if examPaperNum, allErr := ormer.QueryTable(new(models.ExamPaperTest)).Filter("corpid", corpid).Filter("exam_paper_id", examPaperId).Filter("status", 1).All(&examPaperTestsModel, "id", "type", "answer", "score"); allErr != nil || examPaperNum == 0 {
		return nil, allErr
	}

	answersMap := make(map[int]string)
	for _, answer := range answers {
		answerMap, ok := answer.(map[string]interface{})
		if ok {
			answersMap[int(answerMap["id"].(float64))] = answerMap["val"].(string)
		}
	}

	var returnStruct []models.ExamRecordTestDetail
	for _, examPaperTest := range examPaperTestsModel {
		var examRecordTestDetail models.ExamRecordTestDetail
		examRecordTestDetail.Corpid = AttributesStruct.Corpid
		examRecordTestDetail.ExamRecordDetailId = int(examRecordDetailId)
		examRecordTestDetail.ExamPaperTestId = examPaperTest.Id
		examRecordTestDetail.ExamPaperTestType = examPaperTest.Type

		// 未检测到用户答题
		if _, ok := answersMap[examPaperTest.Id]; !ok {
			examRecordTestDetail.Answer = ""
			if isReading > 0 {
				examRecordTestDetail.IsSuccess = 2
				examRecordTestDetail.Score = ""
			} else {
				examRecordTestDetail.IsSuccess = 0
				examRecordTestDetail.Score = "" // 后台阅卷
			}
		} else {
			if isReading > 0 && examPaperTest.Type == 4 { // 需要阅卷
				examRecordTestDetail.Answer = answersMap[examPaperTest.Id]
				examRecordTestDetail.IsSuccess = 2
				examRecordTestDetail.Score = ""
			} else {
				var isSuccess bool

				// eg: 问答题为0分的也是不需要阅卷的
				if examPaperTest.Type == 4 {
					isSuccess = true
					examRecordTestDetail.Answer = answersMap[examPaperTest.Id]
				} else {
					userAnswerString := strings.ToUpper(answersMap[examPaperTest.Id])

					currentAnswerstring := strings.ToUpper(examPaperTest.Answer)
					var currentAnswerSlice []string
					if strings.Index(currentAnswerstring, ",") != -1 {
						currentAnswerSlice = strings.Split(",", currentAnswerstring)
					} else {
						currentAnswerSlice = []string{currentAnswerstring}
					}
					sort.Sort(sort.StringSlice(currentAnswerSlice))

					var currentUserAnswerSlice []string
					if strings.Index(userAnswerString, ",") != -1 {
						currentUserAnswerSlice = strings.Split(",", userAnswerString)
					} else {
						currentUserAnswerSlice = []string{userAnswerString}
					}
					sort.Sort(sort.StringSlice(currentUserAnswerSlice))

					// 题目答案和用户答案比对
					isSuccess = this.StringSliceEqualBCE(currentAnswerSlice, currentUserAnswerSlice)

					examRecordTestDetail.Answer = strings.Join(currentUserAnswerSlice, ",")
				}

				if isSuccess {
					examRecordTestDetail.IsSuccess = 1
					examRecordTestDetail.Score = examPaperTest.Score
				} else {
					examRecordTestDetail.IsSuccess = 0
					examRecordTestDetail.Score = "0"
				}
			}
		}

		returnStruct = append(returnStruct, examRecordTestDetail)
	}

	return returnStruct, nil
}

// 比较两切片是否一样
func (this *ExamController) StringSliceEqualBCE(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	if (a == nil) != (b == nil) {
		return false
	}

	b = b[:len(a)]
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}

// 创建积分记录
func (this *ExamController) addUserScoreLog(examType int8, integral string) (int64, error) {

	now := time.Now()

	dayString := this.FormatTime(int(time.Now().Unix()), "d")
	dayInt64 := this.FormtimeTimeByString(dayString + " 00:00:00")

	var examTypeInfoString string
	if examType == 1 {
		examTypeInfoString = "必考"
	} else {
		examTypeInfoString = "非必考"
	}

	var userScoreModel models.UserScore
	userScoreModel.Corpid = AttributesStruct.Corpid
	userScoreModel.Userid = this.Attributes.Userid
	userScoreModel.Score = integral
	userScoreModel.Type = 7
	userScoreModel.Info = examTypeInfoString + "性质考试通过加" + integral + "分"
	userScoreModel.Day = int(dayInt64)
	userScoreModel.Createtime = int(now.Unix())

	userScoreId, err := orm.NewOrm().Insert(&userScoreModel)

	return userScoreId, err
}
