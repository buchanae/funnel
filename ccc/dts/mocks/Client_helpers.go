package mocks

import dts "github.com/ohsu-comp-bio/funnel/ccc/dts"

func (c *Client) SetFileSites(id string, sites []string) {
  var locs []dts.Location
  for _, site := range sites {
    locs = append(locs, dts.Location{Site: site})
  }
  c.On("GetFile", id).Return(&dts.Entry{
    ID: id,
    Location: locs,
  }, nil)
}
