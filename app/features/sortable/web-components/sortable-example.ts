import { LitElement, html } from 'lit'
import { customElement, property, query } from 'lit/decorators.js'
import Sortable from 'sortablejs'

interface SortableItem {
  name: string
}

@customElement('sortable-example')
export class SortableExample extends LitElement {
  @query('#sortable-container')
  sortContainer!: HTMLElement

  @property({ type: String }) title: string = ''
  @property({ type: String }) value: string = ''
  @property({ type: Array }) items: SortableItem[] = []

  firstUpdated() {
    new Sortable(this.sortContainer, {
      animation: 150,
      ghostClass: 'opacity-25',
      onEnd: (evt) => {
        this.value = `Moved from ${evt.oldIndex} to ${evt.newIndex}`
        this.dispatchEvent(
          new CustomEvent('change', {
            detail: `Moved from ${evt.oldIndex} to ${evt.newIndex}`,
          }),
        )
      },
    })
  }

  protected createRenderRoot() {
    return this
  }

  render() {
    console.log(this)
    return html`
      <div class="sortable-wrapper">
          <div class="sortable-title">${this.title}: <strong>${this.value}</strong></div>
          <div>Open your console to see event results</div>
          <div id="sortable-container" class="sortable-container">
            ${this.items?.length > 0 && this.items.map(
              (item) => html` <div class="sortable-item">${item.name}</div> `,
            )}
        </div>
      </div>
    `
  }
}
