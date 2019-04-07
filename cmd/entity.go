package main

// spatial_ref_sys description goes here ...
type spatial_ref_sys struct {
	srid      int
	auth_name *string
	auth_srid *int
	srtext    *string
	proj4text *string
}

// meta_data description goes here ...
type meta_data struct {
	meta_data_id   *int
	table_name     *string
	style_name     *string
	name           *string
	icon           *string
	pkey           *string
	orderby        *int
	tolerance      *string
	groupname      *string
	legendfilename *string
	order_group    *int
}

// poi_point description goes here ...
type poi_point struct {
	gid      int
	name     *string
	name_en  *string
	name_ru  *string
	man_made *string
	leisure  *string
	amenity  *string
	office   *string
	shop     *string
	tourism  *string
	sport    *string
	osm_type *string
	osm_id   *float64
	geom     *[]byte
}
