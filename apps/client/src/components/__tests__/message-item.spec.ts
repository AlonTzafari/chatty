import { describe, it, expect } from 'vitest'

import { mount } from '@vue/test-utils'
import MessageItem from '../message-item.vue'

describe('MessageItem', () => {
  it('renders properly', () => {
    const wrapper = mount(MessageItem, { 
        props: { 
            message: {
                Content: 'Hello Vitest',
                CreatedAt: new Date().toISOString(),
                Id: 'test-id',
                UserId: 'test-user-id',
                Username: 'vitest-user'
            } 
        } 
    })
    expect(wrapper.text()).toContain('Hello Vitest')
  })
})
