#!/usr/bin/python
import scrapy
from scrapy.selector import Selector

from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker
from sqlalchemy import text

import sys


class RazzSpider(scrapy.Spider):
    db_uri = 'postgresql://localhost/dfs'
    engine = create_engine(db_uri)
    name = 'razzball'
    #start_urls = ['http://razzball.com/dfsbot-fanduel-hit/']
    start_urls = ['http://razzball.com/dfsbot-fanduel-pitch/?uid=1505&token=6eef511269ae262c13541abf93953e8e18bd9b7d15b1b9ea74bf4c86c5cb5818&time=1439075726&loadScript=true',
            'http://razzball.com/dfsbot-fanduel-hit/?uid=1505&token=38c84d51c9e5e3b1a287c39d151ec6c3e6c91881c102655271ec3b67611ea7fd&time=1439085176&loadScript=true']

    def parse(self, response):
        header = response.xpath('//div[@id="content"]/article/header/h1/text()').extract()[0].lower()

        for row in response.xpath('//table[@id="neorazzstatstable"]/tbody/tr'):
            #player = PlayerItem()
            if 'hitting' in header:
                name = row.xpath('td[2]/a/text()').extract()[0]
                team = row.xpath('td[4]/a/text()').extract()[0]
                pos = row.xpath('td[5]/text()').extract()[0]
                vsteam = row.xpath('td[8]/text()').extract()[0]
                pts = row.xpath('td[20]/text()').extract()[0]
                price = row.xpath('td[22]/text()').extract()[0]
            else:
                name = row.xpath('td[2]/a/text()').extract()[0]
                team = row.xpath('td[3]/a/text()').extract()[0]
                pos = 'P'
                vsteam = row.xpath('td[6]/text()').extract()[0]
                pts = row.xpath('td[14]/text()').extract()[0]
                price = row.xpath('td[16]/text()').extract()[0]

            q = ''' INSERT INTO players (name, team, pos, vs_team, points, price) VALUES ('%s', '%s', '%s', '%s', %s, %s) '''
            self.engine.execute(text(q % (name, team, pos, vsteam, pts, price)))
