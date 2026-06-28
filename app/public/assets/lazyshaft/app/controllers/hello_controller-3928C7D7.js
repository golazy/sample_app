import { Controller } from "@hotwired/stimulus"

export default class extends Controller {
  static targets = ["output"]
  static values = {
    speed: { type: Number, default: 90 },
    text: String,
  }

  connect() {
    this.text = this.textValue || this.outputTarget.textContent.trim()
    this.outputTarget.textContent = ""
    this.index = 0
    this.type()
  }

  disconnect() {
    window.clearTimeout(this.timer)
  }

  type() {
    this.outputTarget.textContent = this.text.slice(0, this.index)
    if (this.index >= this.text.length) {
      return
    }
    this.index += 1
    this.timer = window.setTimeout(() => this.type(), this.speedValue)
  }
}
