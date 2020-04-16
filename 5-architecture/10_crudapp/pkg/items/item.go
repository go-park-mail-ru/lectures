package items

type Item struct {
	ID          uint32 `schema:"-"`
	Title       string `schema:"title,required"`
	Description string `schema:"description,required"`
	CreatedBy   uint32 `schema:"-"`
}
