package models

type MenuCursor struct {
	position int
}

func (this *MenuCursor) moveDown(max int, amount int) {
	this.position = min(this.position+amount, max)
}

func (this *MenuCursor) moveUp(amount int) {
	this.position = max(0, this.position-amount)
}
