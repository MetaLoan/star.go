/**
 * Chart Calculation API
 * 本命盘计算接口（无需创建用户）
 */

import { NextRequest, NextResponse } from 'next/server';
import { calculateNatalChart, type BirthData } from '@/lib/astro';

/**
 * POST /api/calc/chart
 * 直接计算本命盘
 */
export async function POST(request: NextRequest) {
  try {
    const body = await request.json();
    
    const { name, date, latitude, longitude, timezone } = body;
    
    if (!date || latitude === undefined || longitude === undefined) {
      return NextResponse.json(
        { error: 'Missing required fields: date, latitude, longitude' },
        { status: 400 }
      );
    }
    
    const birthData: BirthData = {
      name: name || 'Anonymous',
      date: new Date(date),
      latitude,
      longitude,
      timezone: timezone || 'UTC',
    };
    
    const chart = calculateNatalChart(birthData);
    
    // 序列化响应
    return NextResponse.json({
      birthData: {
        name: birthData.name,
        date: birthData.date.toISOString(),
        latitude: birthData.latitude,
        longitude: birthData.longitude,
        timezone: birthData.timezone,
      },
      planets: chart.planets.map(p => ({
        id: p.id,
        name: p.name,
        symbol: p.symbol,
        longitude: p.longitude,
        latitude: p.latitude,
        sign: p.sign,
        signName: p.signName,
        signSymbol: p.signSymbol,
        signDegree: p.signDegree,
        house: p.house,
        retrograde: p.retrograde,
        dignityScore: p.dignityScore,
      })),
      houses: chart.houses.map(h => ({
        house: h.house,
        cusp: h.longitude,
        sign: h.sign,
        signDegree: h.signDegree,
      })),
      ascendant: chart.ascendant,
      midheaven: chart.midheaven,
      aspects: chart.aspects.map(a => ({
        planet1: a.planet1,
        planet2: a.planet2,
        aspectType: a.aspectType,
        exactAngle: a.exactAngle,
        actualAngle: a.actualAngle,
        orb: a.orb,
        strength: a.strength,
        applying: a.applying,
        weight: a.weight,
      })),
      patterns: chart.patterns,
      elementBalance: chart.elementBalance,
      modalityBalance: chart.modalityBalance,
      dominantPlanets: chart.dominantPlanets,
      chartRuler: chart.chartRuler,
    });
  } catch (error) {
    console.error('Chart calculation error:', error);
    return NextResponse.json(
      { error: 'Failed to calculate chart' },
      { status: 500 }
    );
  }
}

