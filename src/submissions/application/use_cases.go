package application

import (
	"mime/multipart"
	"time"

	blocksDefinitions "github.com/UPB-Code-Labs/main-api/src/blocks/domain/definitions"
	"github.com/UPB-Code-Labs/main-api/src/submissions/domain/definitions"
	"github.com/UPB-Code-Labs/main-api/src/submissions/domain/dtos"
	"github.com/UPB-Code-Labs/main-api/src/submissions/domain/entities"
	"github.com/UPB-Code-Labs/main-api/src/submissions/domain/errors"
)

type SubmissionUseCases struct {
	BlocksRepository        blocksDefinitions.BlockRepository
	SubmissionsRepository   definitions.SubmissionsRepository
	SubmissionsQueueManager definitions.SubmissionsQueueManager
}

func (useCases *SubmissionUseCases) CanStudentSubmitToTestBlock(studentUUID string, testBlockUUID string) (bool, error) {
	return useCases.BlocksRepository.CanStudentSubmitToTestBlock(studentUUID, testBlockUUID)
}

func (useCases *SubmissionUseCases) SaveSubmission(dto *dtos.CreateSubmissionDTO) (string, error) {
	// Validate the student can submit to the given test block
	canSubmit, err := useCases.CanStudentSubmitToTestBlock(dto.StudentUUID, dto.TestBlockUUID)
	if err != nil {
		return "", err
	}

	if !canSubmit {
		return "", errors.StudentCannotSubmitToTestBlock{}
	}

	// Check if the student already has a submission for the given test block
	previousStudentSubmission, err := useCases.SubmissionsRepository.GetStudentSubmission(dto.StudentUUID, dto.TestBlockUUID)
	if err != nil {
		return "", err
	}

	if previousStudentSubmission != nil {
		// Check if the previous submission was submitted in the last minute
		parsedSubmittedAt, err := time.Parse(time.RFC3339, previousStudentSubmission.SubmittedAt)
		if err != nil {
			return "", err
		}

		if time.Since(parsedSubmittedAt).Minutes() < 1 {
			return "", errors.StudentHasRecentSubmission{}
		}

		// Check if the previous submission is still pending
		finalStatus := "ready"
		if previousStudentSubmission.Status != finalStatus {
			return "", errors.StudentHasPendingSubmission{}
		}

		// If the student already has a submission, reset its status and overwrite the archive
		err = useCases.resetSubmissionStatus(previousStudentSubmission, dto.SubmissionArchive)
		if err != nil {
			return "", err
		}

		// Submit the work to the submissions queue
		err = useCases.submitWorkToQueue(previousStudentSubmission.UUID)
		if err != nil {
			return "", err
		}

		return previousStudentSubmission.UUID, nil
	} else {
		// If the student doesn't have a submission, create a new one
		submissionUUID, err := useCases.createSubmission(dto)
		if err != nil {
			return "", err
		}

		// Submit the work to the submissions queue
		err = useCases.submitWorkToQueue(submissionUUID)
		if err != nil {
			return "", err
		}

		return submissionUUID, nil
	}
}

func (useCases *SubmissionUseCases) resetSubmissionStatus(previousStudentSubmission *entities.Submission, newArchive *multipart.File) error {
	// Get the UUID of the .zip archive in the static files microservice
	archiveUUID, err := useCases.SubmissionsRepository.GetStudentSubmissionArchiveUUIDFromSubmissionUUID(previousStudentSubmission.UUID)
	if err != nil {
		return err
	}

	// Overwrite the archive in the static files microservice
	err = useCases.SubmissionsRepository.OverwriteSubmissionArchive(newArchive, archiveUUID)
	if err != nil {
		return err
	}

	// Reset the submission status
	err = useCases.SubmissionsRepository.ResetSubmissionStatus(previousStudentSubmission.UUID)
	if err != nil {
		return err
	}

	return nil
}

func (useCases *SubmissionUseCases) createSubmission(dto *dtos.CreateSubmissionDTO) (string, error) {
	// Save the .zip archive in the static files microservice
	archiveUUID, err := useCases.SubmissionsRepository.SaveSubmissionArchive(dto.SubmissionArchive)
	if err != nil {
		return "", err
	}

	dto.SavedArchiveUUID = archiveUUID

	// Save the submission
	submissionUUID, err := useCases.SubmissionsRepository.SaveSubmission(dto)
	if err != nil {
		return "", err
	}

	return submissionUUID, nil
}

func (useCases *SubmissionUseCases) submitWorkToQueue(submissionUUID string) error {
	// Get the submission work
	submissionWork, err := useCases.SubmissionsRepository.GetSubmissionWorkMetadata(submissionUUID)
	if err != nil {
		return err
	}

	// Send the submission work to the submissions queue
	err = useCases.SubmissionsQueueManager.QueueWork(submissionWork)
	if err != nil {
		return err
	}

	return nil
}

func (useCases *SubmissionUseCases) GetSubmissionStatus(studentUUID, testBlockUUID string) (*dtos.SubmissionStatusUpdateDTO, error) {
	// Check if the student could submit to the given test block
	canSubmit, err := useCases.CanStudentSubmitToTestBlock(studentUUID, testBlockUUID)
	if err != nil {
		return nil, err
	}

	if !canSubmit {
		return nil, errors.StudentCannotSubmitToTestBlock{}
	}

	// Get the submission
	submission, err := useCases.SubmissionsRepository.GetStudentSubmission(studentUUID, testBlockUUID)
	if err != nil {
		return nil, err
	}

	if submission == nil {
		return nil, errors.StudentSubmissionNotFound{}
	}

	// Get the submission status
	dto := dtos.SubmissionStatusUpdateDTO{
		SubmissionUUID:   submission.UUID,
		SubmissionStatus: submission.Status,
		TestsPassed:      submission.Passing,
		TestsOutput:      submission.Stdout,
	}

	return &dto, nil
}
