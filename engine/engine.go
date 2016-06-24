package engine

type Context struct {
	Title string
}

func New(title string) *Context {
	ctx := Context{
		Title: title,
	}
	return &ctx
}

func (ctx *Context) Run(scene Scene) error {
	var err error
	running := true

	err = scene.Setup()
	if err != nil {
		return err
	}
	defer scene.Cleanup()

	for _, sub := range scene.SubScenes() {
		err = sub.Setup()
		defer sub.Cleanup()
		if err != nil {
			return err
		}
	}

	for running {
		for _, ent := range scene.Entities() {
			println(ent)
		}

		if scene.ShouldStop() {
			running = false
		}
	}

	return nil
}
