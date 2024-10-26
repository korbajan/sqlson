package dberrors

type NewExecutorError struct {
  dbOpenError error
}

func (e NewExecutorError) Error() string {
  return e.dbOpenError.Error()
}

func CreateNewExecutorError(dbOpenError error) error {
  return NewExecutorError{dbOpenError}
} 
