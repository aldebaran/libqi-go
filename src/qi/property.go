/*
**
** Author(s):
**  - Cedric GESTES <gestes@aldebaran-robotics.com>
**
** Copyright (C) 2013 Aldebaran Robotics
 */

package qi

type Property interface {
	Set(interface{})
	Get() interface{}
	onChanged() Signal
}

type propertyImpl struct {
}

func (p *propertyImpl) Set(v interface{}) {

}

func (p *propertyImpl) Get() interface{} {
	return nil
}

func (p *propertyImpl) onChanged() Signal {
	return NewSignal(func(...interface{}) {})
}

func NewProperty(propertytype interface{}) (Property, error) {
	return Property(&propertyImpl{}), nil
}
