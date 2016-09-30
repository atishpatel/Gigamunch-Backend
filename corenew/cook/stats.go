package cook

// // UpdateNumItemsBy updates a cooks NumItems by the specified count.
// func (c *Client) UpdateNumItemsBy(id string, amount int32) error {
// 	cook, err := get(c.ctx, id)
// 	if err != nil {
// 		return errDatastore.WithError(err).Wrapf("failed to get cook(%s)", id)
// 	}
// 	cook.NumItems += amount
// 	err = put(c.ctx, id, cook)
// 	if err != nil {
// 		return errDatastore.WithError(err).Wrapf("cannot put cook(%s)", id)
// 	}
// 	return nil
// }
