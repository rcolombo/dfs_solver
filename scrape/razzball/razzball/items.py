# -*- coding: utf-8 -*-

# Define here the models for your scraped items
#
# See documentation in:
# http://doc.scrapy.org/en/latest/topics/items.html

import scrapy

class PlayerItem(scrapy.Item):
    name = scrapy.Field()
    team = scrapy.Field()
    pos = scrapy.Field()
    vsteam = scrapy.Field()
    pts = scrapy.Field()
    price = scrapy.Field()
