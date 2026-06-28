import "@hotwired/turbo"
import { Application } from "@hotwired/stimulus"
import HelloController from "/js/controllers/hello_controller.js"

const application = Application.start()
application.register("hello", HelloController)
