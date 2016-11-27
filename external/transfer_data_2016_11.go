package main

import (
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/buslogic"
	"github.com/mijia/sweb/log"
	"time"
)

func TransferDataForNov2016(w *buslogic.Workflow) {
	log.Info("Transfer begin...")

	log.Info("Drop all indexes...")
	if err := w.MongoClient().DropAllIndexes(); err != nil {
		log.Fatalf("fail to drop all indexes: %+v", err)
	}
	log.Info("Done!")
	time.Sleep(5 * time.Second)

	var succ int

	log.Info("1. admin collection...")
	oldAdmins, err := w.MongoClient().GetAllOldAdmins()
	if err != nil {
		log.Fatalf("fail to get old admins: %+v", err)
	}
	succ = 0
	for _, oldAdmin := range oldAdmins {
		newAdmin, err := oldAdmin.ToAdmin()
		if err != nil {
			log.Errorf("fail to convert old admin %+v, err: %+v", oldAdmin, err)
			continue
		}
		if err = w.MongoClient().UpdateAdminWithoutTime(newAdmin); err != nil {
			log.Errorf("fail to update new admin %+v, err: %+v", newAdmin, err)
			continue
		}
		succ++
	}
	log.Infof("Done! %d/%d Succeed", succ, len(oldAdmins))
	time.Sleep(2 * time.Second)

	log.Info("2. teacher collection...")
	oldTeachers, err := w.MongoClient().GetAllOldTeachers()
	if err != nil {
		log.Fatalf("fail to get old teachers: %+v", err)
	}
	succ = 0
	for _, oldTeacher := range oldTeachers {
		newTeacher, err := oldTeacher.ToTeacher()
		if err != nil {
			log.Errorf("fail to convert old teacher %+v, err: %+v", oldTeacher, err)
			continue
		}
		if err = w.MongoClient().UpdateTeacherWithoutTime(newTeacher); err != nil {
			log.Errorf("fail to update new teacher %+v, err: %+v", newTeacher, err)
			continue
		}
		succ++
	}
	log.Infof("Done! %d/%d Succeed", succ, len(oldTeachers))
	time.Sleep(2 * time.Second)

	log.Info("3. student collection...")
	oldStudents, err := w.MongoClient().GetAllOldStudents()
	if err != nil {
		log.Fatalf("fail to get old students: %+v", err)
	}
	succ = 0
	for _, oldStudent := range oldStudents {
		newStudent, err := oldStudent.ToStudent()
		if err != nil {
			log.Errorf("fail to convert old student %+v, err: %+v", oldStudent, err)
			continue
		}
		if err = w.MongoClient().UpdateStudentWithoutTime(newStudent); err != nil {
			log.Errorf("fail to update new student %+v, err: %+v", newStudent, err)
			continue
		}
		succ++
	}
	log.Infof("Done! %d/%d succeed", succ, len(oldStudents))
	time.Sleep(2 * time.Second)

	log.Info("4. archive collection...")
	oldArchives, err := w.MongoClient().GetAllOldArchives()
	if err != nil {
		log.Fatalf("fail to get old archives: %+v", err)
	}
	succ = 0
	for _, oldArchive := range oldArchives {
		newArchive, err := oldArchive.ToArchive()
		if err != nil {
			log.Errorf("fail to convert old archive %+v, err: %+v", oldArchive, err)
			continue
		}
		if err = w.MongoClient().UpdateArchiveWithoutTime(newArchive); err != nil {
			log.Errorf("fail to update new archive %+v, err: %+v", newArchive, err)
			continue
		}
		succ++
	}
	log.Infof("Done! %d/%d succeed", succ, len(oldArchives))
	time.Sleep(2 * time.Second)

	log.Info("5. reservation collection...")
	oldReservations, err := w.MongoClient().GetAllOldReservations()
	if err != nil {
		log.Fatalf("fail to get old reservations: %+v", err)
	}
	succ = 0
	for _, oldReservation := range oldReservations {
		newReservation, err := oldReservation.ToReservation()
		if err != nil {
			log.Errorf("fail to convert old reservation %+v, err: %+v", oldReservation, err)
			continue
		}
		if err = w.MongoClient().UpdateReservationWithoutTime(newReservation); err != nil {
			log.Errorf("fail to update new reservation %+v, err: %+v", newReservation, err)
			continue
		}
		succ++
	}
	log.Infof("Done! %d/%d succeed", succ, len(oldReservations))
	time.Sleep(2 * time.Second)

	log.Info("6. timetable collection...")
	oldTimetables, err := w.MongoClient().GetAllOldTimetables()
	if err != nil {
		log.Fatalf("fail to get old timetables: %+v", err)
	}
	succ = 0
	for _, oldTimetable := range oldTimetables {
		newTimetable, err := oldTimetable.ToTimedReservation()
		if err != nil {
			log.Errorf("fail to convert old timetable %+v, err: %+v", oldTimetable, err)
			continue
		}
		if err = w.MongoClient().UpdateTimedReservationWithoutTime(newTimetable); err != nil {
			log.Errorf("fail to update new timetable %+v, err: %+v", newTimetable, err)
			continue
		}
		succ++
	}
	log.Infof("Done! %d/%d succeed", succ, len(oldTimetables))
	time.Sleep(2 * time.Second)

	log.Info("Ensure all indexes...")
	if err := w.MongoClient().EnsureAllIndexes(); err != nil {
		log.Fatalf("fail to ensure all indexes: %+v", err)
	}
	log.Info("Done!")
	time.Sleep(5 * time.Second)

	log.Info("Transfer Success! Exiting...")
}
